package wallet

import (
	"github.com/stretchr/testify/require"
	"github.com/treeforest/zut.evidence/internal/service/wallet/conf"
	"testing"
)

func TestWallet_Wallet(t *testing.T) {
	uid := int64(100)
	data := "hello world"

	w := New(conf.Default())

	pub, _ := w.GetPubKey(uid)
	if pub == "" {
		var err error
		pub, err = w.GenerateKey(uid)
		require.NoError(t, err)
	}
	t.Logf("public key: %s", pub)

	sig, err := w.Sign(uid, data)
	require.NoError(t, err)
	t.Logf("signatrue: %s", sig)

	ok := w.Verify(sig, pub, data)
	require.Equal(t, true, ok)
}
