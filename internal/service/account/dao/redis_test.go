package dao

import (
	"github.com/stretchr/testify/require"
	"github.com/treeforest/zut.evidence/pkg/jwt"
	redisOp "github.com/treeforest/zut.evidence/pkg/redis"
	"testing"
	"time"
)

func Test_Token(t *testing.T) {
	d := Dao{
		redis:  redisOp.New("localhost:6379", "", 512),
		jwtMgr: jwt.New(jwt.DefaultTokenExpiration * time.Second),
	}
	token, err := d.GenerateToken(10001, 1, "web", []byte{})
	require.NoError(t, err)
	t.Logf("token: %s", token)

	uid, role, platform, extra, err := d.VerifyToken(token)
	require.NoError(t, err)
	t.Logf("uid:%d role:%d platform:%s extra:%v", uid, role, platform, extra)
}
