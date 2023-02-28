package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

const (
	_prefixFileMapping = "fm_"
)

func keyFileMapping(uid int64) string {
	return fmt.Sprintf("%s%d", _prefixFileMapping, uid)
}

// AddFileMapping 添加一条用户所拥有的文件的记录，uid->cid
func (d *Dao) AddFileMapping(uid int64, cid, filename string) error {
	conn := d.redis.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", keyFileMapping(uid), cid, filename)
	return err
}

// ExistFileMapping 判断用户是否拥有cid所对应文件的权限
func (d *Dao) ExistFileMapping(uid int64, cid string) (string, error) {
	conn := d.redis.Get()
	defer conn.Close()
	return redis.String(conn.Do("HGET", keyFileMapping(uid), cid))
}
