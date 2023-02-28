package dao

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/treeforest/zut.evidence/pkg/ipfs"
	"io/ioutil"
	"testing"
)

func Test_Ipfs(t *testing.T) {
	d := &Dao{ipfs: ipfs.New("localhost:5001")}

	cid, err := d.UploadFile(bytes.NewReader([]byte("hello world")))
	require.NoError(t, err)
	require.Equal(t, "Qmf412jQZiuVUtdgnB36FXFX7xg5V6KEbSJ4dpQuhkLyfD", cid)

	rc, err := d.DownloadFile("Qmf412jQZiuVUtdgnB36FXFX7xg5V6KEbSJ4dpQuhkLyfD")
	require.NoError(t, err)

	data, err := ioutil.ReadAll(rc)
	require.NoError(t, err)

	require.NoError(t, err, rc.Close())

	require.Equal(t, "hello world", string(data))
}
