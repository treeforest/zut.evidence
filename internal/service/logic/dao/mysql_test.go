package dao

import (
	"github.com/stretchr/testify/require"
	"github.com/treeforest/zut.evidence/pkg/mysql"
	"testing"
)

func getTestDao() *Dao {
	return &Dao{mysql: mysql.Connect("user", "123456", "localhost:3306", "zut")}
}

func TestDao_CreateTables(t *testing.T) {
	d := getTestDao()
	d.createTables()
}

func TestDao_CreateApply(t *testing.T) {
	d := getTestDao()

	for i := 0; i < 10; i++ {
		_, err := d.CreateApply(101, 202, "sender_did", "recipient_did",
			Education, "this is a reason", "this is some cids")
		require.NoError(t, err)
	}
}

func TestDao_GetApplyDoing(t *testing.T) {
	d := getTestDao()

	contents, err := d.GetApplyDoing(10004)
	require.NoError(t, err)

	for i, content := range contents {
		t.Logf("[%d] %v", i, content)
	}
}

func TestDao_GetApplyDoingCount(t *testing.T) {
	d := getTestDao()

	count, err := d.GetApplyDoingCount(101)
	require.NoError(t, err)

	t.Logf("count: %d", count)
}

func TestDao_GetAuditDoing(t *testing.T) {
	d := getTestDao()

	contents, err := d.GetAuditDoing(10006)
	require.NoError(t, err)

	for i, content := range contents {
		t.Logf("[%d] %v", i, content)
	}
}
