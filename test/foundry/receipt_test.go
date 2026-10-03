package foundry_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-contracts/test/foundry"
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

	t.Run("happy path - missing blob gas price", func(t *testing.T) {
		// given
		data := map[string]any{}
		require.NoError(t, json.Unmarshal(receiptJSON, &data))
		delete(data, "blobGasPrice")
		want, err := json.Marshal(data)
		require.NoError(t, err)

		// when
		receipt := &foundry.Receipt{}
		err = receipt.UnmarshalJSON(want)

		// then
		require.NoError(t, err)
		require.Zero(t, receipt.BlobGasPrice)
	})
}
