package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
)

const (
	Type = "Secp256k1"
)

var (
	curve = elliptic.P256()
)

// GenerateKey 生成私钥
func GenerateKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return key, nil
}

// Sign 签名
func Sign(key *ecdsa.PrivateKey, hash [32]byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, key, hash[:])
}

// Verify 验签
func Verify(pub *ecdsa.PublicKey, signature []byte, hash [32]byte) bool {
	return ecdsa.VerifyASN1(pub, hash[:], signature)
}

// MarshalKey 序列化私钥
func MarshalKey(key *ecdsa.PrivateKey) []byte {
	keyAsn1Data, _ := x509.MarshalECPrivateKey(key)
	return keyAsn1Data
}

// UnmarshalKey 反序列化私钥
func UnmarshalKey(data []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(data)
}

// MarshalKeyPem 将私钥转换为pem格式
func MarshalKeyPem(key *ecdsa.PrivateKey) []byte {
	keyAsn1Data := MarshalKey(key)
	block := &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: keyAsn1Data,
	}
	return pem.EncodeToMemory(block)
}

// UnmarshalKeyPem 从pem格式中解析出私钥
func UnmarshalKeyPem(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("pem format error")
	}
	return UnmarshalKey(block.Bytes)
}

// MarshalPubKey 序列化公钥
func MarshalPubKey(pub *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(curve, pub.X, pub.Y)
}

// UnmarshalPubKey 反序列化公钥
func UnmarshalPubKey(data []byte) (*ecdsa.PublicKey, error) {
	pub := ecdsa.PublicKey{Curve: curve}
	pub.X, pub.Y = elliptic.Unmarshal(pub.Curve, data)
	if pub.X == nil || pub.Y == nil {
		return nil, errors.New("unmarshal failed")
	}
	return &pub, nil
}

// MarshalPubKeyPem 将公钥转换为pem格式
func MarshalPubKeyPem(pub *ecdsa.PublicKey) []byte {
	block := &pem.Block{
		Type:  "ECDSA PUBLIC KEY",
		Bytes: MarshalPubKey(pub),
	}
	return pem.EncodeToMemory(block)
}

// UnmarshalPubKeyPem 从pem格式中解析出公钥
func UnmarshalPubKeyPem(data []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("pem format error")
	}
	return UnmarshalPubKey(block.Bytes)
}
