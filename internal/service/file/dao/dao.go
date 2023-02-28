package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/treeforest/zut.evidence/internal/service/file/conf"
	"github.com/treeforest/zut.evidence/pkg/ipfs"
	redisHelper "github.com/treeforest/zut.evidence/pkg/redis"
)

type Dao struct {
	ipfs  *ipfs.Ipfs
	redis *redis.Pool
}

func New(c *conf.Config) *Dao {
	d := &Dao{
		ipfs:  ipfs.New(c.Ipfs.Url),
		redis: redisHelper.New(c.Redis.Addr, c.Redis.Password, c.Redis.Maxidle),
	}
	return d
}

func (d *Dao) Close() {
	_ = d.redis.Close()
}
