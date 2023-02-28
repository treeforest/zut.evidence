package dao

import (
	"github.com/treeforest/zut.evidence/pkg/mysql"
	"testing"
)

func Test_Mysql(t *testing.T) {
	d := &Dao{mysql: mysql.Connect("user", "123456", "localhost:3306", "zut")}

	d.createTables()

	err := d.Register("tony", "15690087654", "123@qq.com", "1234", 1)
	t.Logf("err:%v", err)

	uid, err := d.Login("tony", "1234", 1)
	t.Logf("uid:%d err:%v", uid, err)
}
