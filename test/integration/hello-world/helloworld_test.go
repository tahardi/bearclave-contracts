package helloworld_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/contracts/bindings"
	"github.com/tahardi/bearchain/test/foundry"
)

func TestHelloWorld(t *testing.T) {
	// given
	broadcastDir := "../../../contracts/broadcast"
	scriptDir := "../../../contracts/scripts"

	anvil, err := foundry.NewAnvil(broadcastDir, scriptDir)
	require.NoError(t, err)

	err = anvil.Start()
	require.NoError(t, err)
	defer anvil.Stop()

	owner := anvil.Account(0)
	contractName := "HelloWorld"
	contractAddress, err := anvil.DeployContract(contractName, owner)
	require.NoError(t, err)

	client, err := ethclient.Dial(anvil.URL())
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
