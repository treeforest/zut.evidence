package dao

import (
	"github.com/treeforest/zut.evidence/internal/service/logic/conf"
	"github.com/treeforest/zut.evidence/pkg/mysql"

	"gorm.io/gorm"
)

type Dao struct {
	mysql *gorm.DB
}

func New(c *conf.Config) *Dao {
	d := &Dao{
		mysql: mysql.Connect(c.Mysql.User, c.Mysql.Pass, c.Mysql.Addr, c.Mysql.Database),
	}
	d.createTables()
	return d
}

func (d *Dao) Close() {
}
