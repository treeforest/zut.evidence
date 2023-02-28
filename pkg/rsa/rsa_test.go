package rsa

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RSA(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	data := []byte("hello world")

	// 签名验签
	sig, err := Sign(key, data)
	require.NoError(t, err)
	err = Verify(&key.PublicKey, data, sig)
	require.NoError(t, err)

	// 加解密
	cipher, err := Encrypt(&key.PublicKey, data)
	require.NoError(t, err)
	b, err := Decrypt(key, cipher)
	require.NoError(t, err)
	require.Equal(t, data, b)
}
