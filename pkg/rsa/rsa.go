package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
)

// GenerateKey 密钥生成
func GenerateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// Sign 签名
func Sign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
}

// Verify 验签
func Verify(pub *rsa.PublicKey, data, sig []byte) error {
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig)
}

// Encrypt 加密
func Encrypt(pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
}

// Decrypt 解密
func Decrypt(key *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext)
}

// MarshalKey 序列化私钥
func MarshalKey(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

// UnmarshalKey 反序列化私钥
func UnmarshalKey(data []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(data)
}

// MarshalKeyPem 将私钥转换为pem格式
func MarshalKeyPem(key *rsa.PrivateKey) []byte {
	keyAsn1Data := MarshalKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyAsn1Data,
	}
	return pem.EncodeToMemory(block)
}

// UnmarshalKeyPem 从pem格式中解析出私钥
func UnmarshalKeyPem(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("pem format error")
	}
	return UnmarshalKey(block.Bytes)
}

// MarshalPubKey 序列化公钥
func MarshalPubKey(pub *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pub)
}

// UnmarshalPubKey 反序列化公钥
func UnmarshalPubKey(data []byte) (*rsa.PublicKey, error) {
	return x509.ParsePKCS1PublicKey(data)
}

// MarshalPubKeyPem 将公钥转换为pem格式
func MarshalPubKeyPem(pub *rsa.PublicKey) []byte {
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: MarshalPubKey(pub),
	}
	return pem.EncodeToMemory(block)
}

// UnmarshalPubKeyPem 从pem格式中解析出公钥
func UnmarshalPubKeyPem(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("pem format error")
	}
	return UnmarshalPubKey(block.Bytes)
}
