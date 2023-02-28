package did

import (
	"fmt"
	"time"
)

const (
	version       = 1
	SHA256WithRSA = "SHA256WithRSA"
)

var (
	defaultContext = []string{
		"https://www.w3.org/ns/did/v1",
	}
)

// Proof 证明
type Proof struct {
	Type      string `json:"type"`      // 签名算法
	Creator   string `json:"creator"`   // 签名者的公钥id
	Signature string `json:"signature"` // 签名
}

// PublicKey 公钥
type PublicKey struct {
	Id           string `json:"id"`           // 公钥id，格式为：did+#key-1
	Type         string `json:"type"`         // 密钥生成算法
	PublicKeyHex string `json:"publicKeyHex"` // 公钥的十六进制表示
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GeneratePublicKeyID(did string, id int) string {
	return fmt.Sprintf("%s#keys-%d", did, id)
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
