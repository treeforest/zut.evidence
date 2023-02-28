package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
	redisHelper "github.com/treeforest/zut.evidence/pkg/redis"
	"testing"
)

func Test_FileMapping(t *testing.T) {
	d := Dao{
		redis: redisHelper.New("localhost:6379", "", 512),
	}
	defer d.redis.Close()

	require.NoError(t, d.AddFileMapping(101, "123", "aa.txt"))
	require.NoError(t, d.AddFileMapping(102, "456", "bb.txt"))
	require.NoError(t, d.AddFileMapping(103, "789", "cc.txt"))

	filename, err := d.ExistFileMapping(101, "123")
	require.NoError(t, err)
	require.Equal(t, "aa.txt", filename)

	filename, err = d.ExistFileMapping(101, "456")
	require.Equal(t, redis.ErrNil, err)
	require.Equal(t, "", filename)
}
