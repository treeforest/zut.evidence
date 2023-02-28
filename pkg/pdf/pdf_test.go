package pdf

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"github.com/treeforest/zut.evidence/pkg/did"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func TestGenEvidencePdf(t *testing.T) {
	UtilPath = path.Join(".")

	issuerKey, err := rsaHelper.GenerateKey()
	require.NoError(t, err)
	issuerPublicKeyHex := hex.EncodeToString(rsaHelper.MarshalPubKey(&issuerKey.PublicKey))
	issuerDocument := did.NewDocument(issuerPublicKeyHex, issuerKey)

	userKey, err := rsaHelper.GenerateKey()
	require.NoError(t, err)
	userPublicKeyHex := hex.EncodeToString(rsaHelper.MarshalPubKey(&userKey.PublicKey))
	userDocument := did.NewDocument(userPublicKeyHex, issuerKey)

	proofClaim, err := did.IssueProofClaim(time.Hour*24, "实名认证", issuerDocument.Id,
		issuerDocument.PublicKey[0].Id, userDocument.Id, "张三实名认证",
		"张三居住在马路边的犄角旮瘩张三居住在马路边的犄角旮瘩张三居住在马路边的犄角旮瘩", issuerKey)
	require.NoError(t, err)

	data := GenEvidencePdf(&PdfInfo{
		ProofClaim:  proofClaim,
		DownloadUrl: "http://www.baidu.com",
	})

	err = ioutil.WriteFile("test.pdf", data, 777)
	require.NoError(t, err)
}
