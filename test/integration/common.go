package integration

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/test/foundry"
)

const (
	ContractDir  = "../../../contracts"
	BroadcastDir = ContractDir + "/broadcast"
	ScriptDir    = ContractDir + "/scripts"
)

func AssertAddressesEqual(
	t *testing.T,
	address1 common.Address,
	address2 common.Address,
) {
	t.Helper()
	assert.Equal(t, 0, address1.Cmp(address2))
}

func StartAnvil(
	t *testing.T,
	silent bool,
) (*foundry.Anvil, func()) {
	t.Helper()
	anvil, err := foundry.NewAnvil(BroadcastDir, ScriptDir)
	require.NoError(t, err)

	err = anvil.Start(t.Context(), silent)
	require.NoError(t, err)

	stop := func() { _ = anvil.Stop() }
	return anvil, stop
}
