package jwt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	jwtMgr                         *JWTManager
	_uid, _role, _platform, _extra = int64(10001), 1, "web", []byte("extra message")
	_fakeToken                     string
)

func init() {
	jwtMgr = New(time.Second)
	_fakeToken, _ = jwtMgr.Generate(_uid, _role, _platform, _extra)
}

func Test_JWTManager(t *testing.T) {
	// 生成token
	token, err := jwtMgr.Generate(_uid, _role, _platform, _extra)
	require.NoError(t, err)

	// 正确解析
	uid, role, platform, extra, err := jwtMgr.Verify(token)
	require.NoError(t, err)
	require.Equal(t, _uid, uid)
	require.Equal(t, _role, role)
	require.Equal(t, _platform, platform)
	require.Equal(t, _extra, extra)

	time.Sleep(time.Second * 2)

	// 过期
	_, _, _, _, err = jwtMgr.Verify(token)
	t.Log(err)
}

func Benchmark_Generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = jwtMgr.Generate(_uid, _role, _platform, _extra)
	}
}

func Benchmark_Verify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _, _ = jwtMgr.Verify(_fakeToken)
	}
}
