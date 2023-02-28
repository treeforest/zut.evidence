package did

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"time"
)

// ProofClaim 可验证凭证（可理解为数字证书）
type ProofClaim struct {
	Context           []string          `json:"@context"`
	Id                string            `json:"id"`             // claim唯一id
	Type              string            `json:"type"`           // 类型，固定为ProofClaim
	Issuer            string            `json:"issuer"`         // 颁发者did
	IssuanceDate      string            `json:"issuanceDate"`   // 颁发日期
	ExpirationDate    string            `json:"expirationDate"` // 过期时间
	CredentialSubject CredentialSubject `json:"credentialSubject"`
	Proof             Proof             `json:"proof"`
}

// CredentialSubject 凭证主体
type CredentialSubject struct {
	Id               string `json:"id"`               // 拥有者的did
	ShortDescription string `json:"shortDescription"` // 对该凭证的简短描述
	LongDescription  string `json:"longDescription"`  // 对该凭证的详细描述
	Type             string `json:"type"`             // 凭证类型（实名认证、学历认证、工作经历认证等等）
}

func (pc *ProofClaim) Marshal() []byte {
	b, _ := json.Marshal(pc)
	return b
}

func (pc *ProofClaim) Unmarshal(b []byte) error {
	return json.Unmarshal(b, pc)
}

func BaseProofClaim(expirationData string, claimType string, issuerDID, issuerPubKeyID,
	holderDID, shortDesc, longDesc string) *ProofClaim {

	issuanceTime := time.Now()
	return &ProofClaim{
		Context:        defaultContext,
		Id:             uuid.New().String(),
		Type:           "ProofClaim",
		Issuer:         issuerDID,
		IssuanceDate:   TimeFormat(issuanceTime),
		ExpirationDate: expirationData, // TimeFormat(issuanceTime.Add(expirationDuration)),
		CredentialSubject: CredentialSubject{
			Id:               holderDID,
			ShortDescription: shortDesc,
			LongDescription:  longDesc,
			Type:             claimType,
		},
		Proof: Proof{
			Type:      SHA256WithRSA,
			Creator:   issuerPubKeyID,
			Signature: "",
		},
	}
}

// IssueProofClaim 颁发可验证声明
func IssueProofClaim(expirationDuration time.Duration, claimType string, issuerDID, issuerPubKeyID,
	holderDID, shortDesc, longDesc string, key *rsa.PrivateKey) (*ProofClaim, error) {
	issuanceTime := time.Now()
	pc := &ProofClaim{
		Context:        defaultContext,
		Id:             uuid.New().String(),
		Type:           "ProofClaim",
		Issuer:         issuerDID,
		IssuanceDate:   TimeFormat(issuanceTime),
		ExpirationDate: TimeFormat(issuanceTime.Add(expirationDuration)),
		CredentialSubject: CredentialSubject{
			Id:               holderDID,
			ShortDescription: shortDesc,
			LongDescription:  longDesc,
			Type:             claimType,
		},
		Proof: Proof{
			Type:      SHA256WithRSA,
			Creator:   issuerPubKeyID,
			Signature: "",
		},
	}

	data := pc.Marshal()
	signature, err := rsaHelper.Sign(key, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pc.Proof.Signature = hex.EncodeToString(signature)

	return pc, nil
}
