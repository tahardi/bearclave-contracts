package helloworld_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-contracts/contracts/bindings"
	"github.com/tahardi/bearclave-contracts/test/integration"
)

const (
	ContractName = "HelloWorld"
)

func TestHelloWorld(t *testing.T) {
	// given
	f, stop := integration.StartFoundry(t, true)
	defer stop()

	owner := f.Anvil().Account(0)
	contractAddress, err := f.Forge().DeployContract(t.Context(), ContractName, owner)
	require.NoError(t, err)

	client, err := ethclient.Dial(f.Anvil().URL())
	require.NoError(t, err)

	want := "Hello, World!"
	hwContract, err := bindings.NewHelloWorld(*contractAddress, client)
	require.NoError(t, err)

	// when
	greeting, err := hwContract.Greet(nil)

	// then
	require.NoError(t, err)
	assert.Equal(t, want, greeting)
}
