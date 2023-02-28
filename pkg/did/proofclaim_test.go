package did

import (
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/require"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"testing"
	"time"
)

func TestIssueProofClaim(t *testing.T) {
	issuerKey, err := rsaHelper.GenerateKey()
	require.NoError(t, err)
	issuerPublicKeyHex := hex.EncodeToString(rsaHelper.MarshalPubKey(&issuerKey.PublicKey))
	issuerDocument := NewDocument(issuerPublicKeyHex, issuerKey)

	userKey, err := rsaHelper.GenerateKey()
	require.NoError(t, err)
	userPublicKeyHex := hex.EncodeToString(rsaHelper.MarshalPubKey(&userKey.PublicKey))
	userDocument := NewDocument(userPublicKeyHex, issuerKey)

	proofClaim, err := IssueProofClaim(time.Hour*24, "实名认证", issuerDocument.Id,
		issuerDocument.PublicKey[0].Id, userDocument.Id, "张三实名认证", "张三居住在**马路边的犄角旮瘩", issuerKey)
	require.NoError(t, err)

	b, _ := json.MarshalIndent(proofClaim, "", "\t")
	t.Logf("ProofClaim:\n%s", b)
}
