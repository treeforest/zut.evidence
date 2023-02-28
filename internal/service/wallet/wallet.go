package wallet

import (
	"crypto/rsa"
	"encoding/hex"
	log "github.com/treeforest/logger"
	"github.com/treeforest/zut.evidence/internal/service/wallet/conf"
	"github.com/treeforest/zut.evidence/internal/service/wallet/dao"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Wallet struct {
	dao *dao.Dao
}

func New(c *conf.Config) *Wallet {
	return &Wallet{dao: dao.New(c)}
}

func (w *Wallet) Close() error {
	w.dao.Close()
	return nil
}

// GenerateKey 生成公私钥
func (w *Wallet) GenerateKey(uid int64) (string, error) {
	exist, err := w.dao.ExistKey(uid)
	if err != nil {
		return "", err
	}
	if exist {
		return "", status.Error(codes.InvalidArgument, "密钥已存在")
	}

	key, err := rsaHelper.GenerateKey()
	if err != nil {
		return "", status.Errorf(codes.Internal, "创建私钥失败: %v", err)
	}
	err = w.dao.AddKey(uid, hex.EncodeToString(rsaHelper.MarshalKey(key)))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rsaHelper.MarshalPubKey(&key.PublicKey)), nil
}

// DownloadKey 下载公私钥
func (w *Wallet) DownloadKey(uid int64) ([]byte, []byte, error) {
	key, err := w.getKey(uid)
	if err != nil {
		return nil, nil, err
	}
	return rsaHelper.MarshalKeyPem(key), rsaHelper.MarshalPubKeyPem(&key.PublicKey), nil
}

// GetPrivateKey 获取私钥
func (w *Wallet) GetPrivateKey(uid int64) ([]byte, error) {
	key, err := w.getKey(uid)
	if err != nil {
		return nil, err
	}
	return rsaHelper.MarshalKey(key), nil
}

// GetPubKey 获取公钥
func (w *Wallet) GetPubKey(uid int64) (string, error) {
	key, err := w.getKey(uid)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rsaHelper.MarshalPubKey(&key.PublicKey)), nil
}

// Sign 签名
func (w *Wallet) Sign(uid int64, data string) (string, error) {
	key, err := w.getKey(uid)
	if err != nil {
		return "", err
	}
	sig, err := rsaHelper.Sign(key, []byte(data))
	if err != nil {
		return "", status.Errorf(codes.Internal, "签名失败: %v", err)
	}
	return hex.EncodeToString(sig), nil
}

// Verify 验签
func (w *Wallet) Verify(hexSig, hexPub, data string) bool {
	sig, err := hex.DecodeString(hexSig)
	if err != nil {
		log.Debug("failed to decode signature: ", err)
		return false
	}
	pubBytes, err := hex.DecodeString(hexPub)
	if err != nil {
		log.Debug("failed to decode public: ", err)
		return false
	}
	pub, err := rsaHelper.UnmarshalPubKey(pubBytes)
	if err != nil {
		log.Debug("failed to unmarshal public:", err)
		return false
	}

	if err = rsaHelper.Verify(pub, []byte(data), sig); err != nil {
		log.Debug("failed to verify: ", err)
		return false
	}
	return true
}

// Encrypt 加密
func (w *Wallet) Encrypt(uid int64, data string) (string, error) {
	key, err := w.getKey(uid)
	if err != nil {
		return "", err
	}
	ciphertext, err := rsaHelper.Encrypt(&key.PublicKey, []byte(data))
	if err != nil {
		return "", status.Errorf(codes.Internal, "加密失败: %v", err)
	}
	return hex.EncodeToString(ciphertext), nil
}

// EncryptByPubKey 加密
func (w *Wallet) EncryptByPubKey(data, hexPub string) (string, error) {
	pubBytes, err := hex.DecodeString(hexPub)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "解码十六进制公钥失败: %v", err.Error())
	}
	pub, err := rsaHelper.UnmarshalPubKey(pubBytes)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "反序列化公钥失败: %v", err.Error())
	}
	ciphertext, err := rsaHelper.Encrypt(pub, []byte(data))
	if err != nil {
		return "", status.Errorf(codes.Internal, "加密失败: %v", err.Error())
	}
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt 解密
func (w *Wallet) Decrypt(uid int64, ciphertext string) (string, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "无效的加密文本: %v", err.Error())
	}
	key, err := w.getKey(uid)
	if err != nil {
		return "", err
	}
	data, err := rsaHelper.Decrypt(key, ciphertextBytes)
	if err != nil {
		return "", status.Errorf(codes.Internal, "解密失败: %v", err)
	}
	return string(data), nil
}

func (w *Wallet) getKey(uid int64) (*rsa.PrivateKey, error) {
	keyStr, err := w.dao.GetKey(uid)
	if err != nil {
		return nil, err
	}
	keyBytes, _ := hex.DecodeString(keyStr)
	key, err := rsaHelper.UnmarshalKey(keyBytes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "私钥反序列化失败: %v", err)
	}
	return key, err
}
