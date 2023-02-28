package did

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/treeforest/base58check"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"golang.org/x/crypto/ripemd160"
)

// Document did document 的定义
type Document struct {
	Context        []string    `json:"@context"`
	Id             string      `json:"id"`
	Version        int         `json:"version"`
	Created        string      `json:"created"`
	Updated        string      `json:"updated"`
	PublicKey      []PublicKey `json:"publicKey"`
	Authentication []string    `json:"authentication"`
	Proof          Proof       `json:"proof"`
}

func (doc *Document) Marshal() []byte {
	b, _ := json.Marshal(doc)
	return b
}

func (doc *Document) Unmarshal(b []byte) error {
	return json.Unmarshal(b, doc)
}

func (doc *Document) GenerateDID() string {
	return doc.Id
}

func NewDocument(publicKeyHex string, key *rsa.PrivateKey) *Document {
	curTime := CurrentTime()
	baseDocument := &Document{
		Context: defaultContext,
		Id:      "",
		Version: version,
		Created: curTime,
		Updated: curTime,
		PublicKey: []PublicKey{
			{
				Id:           "#keys-1",
				Type:         SHA256WithRSA,
				PublicKeyHex: publicKeyHex,
			},
		},
		Authentication: []string{"#keys-1"},
		Proof: Proof{
			Type:      SHA256WithRSA,
			Creator:   "#keys-1",
			Signature: "",
		},
	}
	did := GenerateDIDWithBaseDocument(baseDocument)
	publicKeyID := GeneratePublicKeyID(did, 1)

	baseDocument.Id = did
	baseDocument.PublicKey[0].Id = publicKeyID
	baseDocument.Authentication[0] = publicKeyID
	baseDocument.Proof.Creator = publicKeyID

	signature, _ := rsaHelper.Sign(key, baseDocument.Marshal())
	baseDocument.Proof.Signature = hex.EncodeToString(signature)

	return baseDocument
}

func GenerateDIDWithBaseDocument(baseDocument *Document) string {
	baseDocumentBytes := baseDocument.Marshal()
	hash160 := calcHash160(baseDocumentBytes)
	return fmt.Sprintf("did:zut:%s", base58check.Encode(hash160))
}

func CreateDocumentWithPrivateKey(key *rsa.PrivateKey) *Document {
	publicKeyHex := hex.EncodeToString(rsaHelper.MarshalPubKey(&key.PublicKey))
	return NewDocument(publicKeyHex, key)
}

func calcHash160(data []byte) []byte {
	hash := sha256.Sum256(data)
	r := ripemd160.New()
	r.Write(hash[:])
	hash160 := r.Sum(nil)
	return hash160[:]
}
