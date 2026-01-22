package foundry_test

import (
	_ "embed"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/test/foundry"
)

const (
	contractName    = "BearCoin"
	contractAddress = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
)

//go:embed testdata/broadcast.json
var broadcastJSON []byte

func TestBroadcast_GetContractAddress(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		broadcast := &foundry.Broadcast{}
		require.NoError(t, json.Unmarshal(broadcastJSON, broadcast))

		// when
		got, err := broadcast.GetContractAddress(contractName)

		//
		require.NoError(t, err)
		require.Equal(t, contractAddress, strings.ToLower(got.Hex()))
	})

	t.Run("error - contract not found", func(t *testing.T) {
		// given
		broadcast := &foundry.Broadcast{}
		require.NoError(t, json.Unmarshal(broadcastJSON, broadcast))

		// when
		_, err := broadcast.GetContractAddress("HelloWorld")

		// then
		require.ErrorIs(t, err, foundry.ErrContractNotFound)
	})
}

func TestBroadcast_JSON(t *testing.T) {
	t.Run("happy path - round trip", func(t *testing.T) {
		// given
		want := broadcastJSON

		// when
		broadcast := &foundry.Broadcast{}
		err := json.Unmarshal(want, broadcast)
		require.NoError(t, err)

		got, err := json.Marshal(broadcast)

		// then
		require.NoError(t, err)
		require.JSONEq(t, string(want), string(got))
	})
}
