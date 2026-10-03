package integration

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-foundry/foundry"
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

func StartFoundry(
	t *testing.T,
	silent bool,
) (*foundry.Foundry, func()) {
	t.Helper()
	f, err := foundry.NewFoundry(t.Context(), silent, BroadcastDir, ScriptDir)
	require.NoError(t, err)
	return f, f.Stop
}
