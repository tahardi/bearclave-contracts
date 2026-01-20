package helloworld_test

import (
	"net"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/contracts/bindings"
)

func TestHelloWorld(t *testing.T) {
	// given
	host, ok := os.LookupEnv("HOST")
	require.True(t, ok)
	port, ok := os.LookupEnv("PORT")
	require.True(t, ok)
	contractAddress, ok := os.LookupEnv("CONTRACT_ADDRESS")
	require.True(t, ok)

	url := "http://" + net.JoinHostPort(host, port)
	client, err := ethclient.Dial(url)
	require.NoError(t, err)

	want := "Hello, World!"
	hwContract, err := bindings.NewHelloWorld(common.HexToAddress(contractAddress), client)
	require.NoError(t, err)

	// when
	greeting, err := hwContract.Greet(nil)

	// then
	require.NoError(t, err)
	assert.Equal(t, want, greeting)
}
