package ecdsa

import (
	"crypto/sha256"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)

	hash := sha256.Sum256([]byte("hello world"))

	// 签名
	sig, err := Sign(key, hash)
	require.NoError(t, err)

	// 验签
	ok := Verify(&key.PublicKey, sig, hash)
	require.Equal(t, true, ok)

	// 序列化与反序列化私钥
	data := MarshalKeyPem(key)
	t.Logf("private key:\n%s", data)

	key2, err := UnmarshalKeyPem(data)
	require.NoError(t, err)
	require.Equal(t, true, key2.Equal(key))

	// 序列化与反序列化公钥
	data = MarshalPubKeyPem(&key.PublicKey)
	t.Logf("public key:\n%s", data)

	pub, err := UnmarshalPubKeyPem(data)
	require.NoError(t, err)
	require.Equal(t, true, pub.Equal(&key.PublicKey))
}
