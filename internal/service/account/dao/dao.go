package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/treeforest/zut.evidence/internal/service/account/conf"
	jwt "github.com/treeforest/zut.evidence/pkg/jwt"
	"github.com/treeforest/zut.evidence/pkg/mysql"
	redisOp "github.com/treeforest/zut.evidence/pkg/redis"
	"gorm.io/gorm"
)

type Dao struct {
	mysql  *gorm.DB
	redis  *redis.Pool
	jwtMgr *jwt.JWTManager
}

func New(c *conf.Config, jwtMgr *jwt.JWTManager) *Dao {
	d := &Dao{
		mysql:  mysql.Connect(c.Mysql.User, c.Mysql.Pass, c.Mysql.Addr, c.Mysql.Database),
		redis:  redisOp.New(c.Redis.Addr, c.Redis.Password, c.Redis.Maxidle),
		jwtMgr: jwtMgr,
	}
	d.createTables()
	return d
}

func (d *Dao) Close() {
	_ = d.redis.Close()
}
