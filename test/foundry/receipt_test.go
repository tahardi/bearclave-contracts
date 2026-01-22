package foundry_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/test/foundry"
)

//go:embed testdata/receipt.json
var receiptJSON []byte

func TestReceipt_JSON(t *testing.T) {
	t.Run("happy path - round trip", func(t *testing.T) {
		// given
		want := receiptJSON

		// when
		receipt := &foundry.Receipt{}
		err := receipt.UnmarshalJSON(want)
		require.NoError(t, err)

		got, err := receipt.MarshalJSON()

		// then
		require.NoError(t, err)
		require.JSONEq(t, string(want), string(got))
	})
}
