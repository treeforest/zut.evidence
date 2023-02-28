package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/treeforest/zut.evidence/pkg/jwt"
)

const (
	_prefixToken = "token_"
)

func keyToken(uid int64, platform string) string {
	return fmt.Sprintf("%s%d_%s", _prefixToken, uid, platform)
}

// GenerateToken 根据条件生成token
func (d *Dao) GenerateToken(uid int64, role int, platform string, extra []byte) (string, error) {
	token, err := d.jwtMgr.Generate(uid, role, platform, extra)
	if err != nil {
		return "", errors.WithStack(err)
	}
	conn := d.redis.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", keyToken(uid, platform), jwt.DefaultTokenExpiration, token)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return token, nil
}

// VerifyToken 验证token，若验证成功，则一并返回token中所包含的信息
func (d *Dao) VerifyToken(token string) (uid int64, role int, platform string, extra []byte, err error) {
	uid, role, platform, extra, err = d.jwtMgr.Verify(token)
	if err != nil {
		return 0, 0, "", nil, errors.WithStack(err)
	}
	conn := d.redis.Get()
	defer conn.Close()
	cacheToken, err := redis.String(conn.Do("GET", keyToken(uid, platform)))
	if err != nil {
		return 0, 0, "", nil, errors.WithStack(err)
	}
	if cacheToken != token {
		return 0, 0, "", nil, errors.New("invalid token")
	}
	return uid, role, platform, extra, nil
}
