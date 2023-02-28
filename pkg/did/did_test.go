package did

import (
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/require"
	rsaHelper "github.com/treeforest/zut.evidence/pkg/rsa"
	"testing"
)

func TestNewDocument(t *testing.T) {
	key, err := rsaHelper.GenerateKey()
	require.NoError(t, err)
	document := NewDocument(hex.EncodeToString(rsaHelper.MarshalPubKey(&key.PublicKey)), key)
	b, err := json.MarshalIndent(document, "", "\t")
	require.NoError(t, err)
	t.Logf("DID Document:\n%s", string(b))
}
