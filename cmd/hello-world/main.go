package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tahardi/bearchain/contracts/bindings"
)

const (
	EmptyAddress = ""
	DefaultHost = "127.0.0.1"
	DefaultPort = 8545
)

var (
	address string
	host string
	port int
)

func main() {
	flag.StringVar(
		&address,
		"address",
		EmptyAddress,
		"The address of the contract to interact with",
		)
	flag.StringVar(
		&host,
		"host",
		DefaultHost,
		"The hostname of the eth client to connect to (default: 127.0.0.1)",
	)
	flag.IntVar(
		&port,
		"port",
		DefaultPort,
		"The port of the eth client to connect to (default: 8545)",
	)
	flag.Parse()

	if address == EmptyAddress {
		log.Fatal("missing contract address")
	}

	url := "http://" + net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("connecting to eth client '%s': %v", url, err)
	}

	contractAddr := common.HexToAddress(address)
	hwContract, err := bindings.NewHelloWorld(contractAddr, client)
	if err != nil {
		log.Fatalf("creating hello world contract '%s': %v", address, err)
	}

	greeting, err := hwContract.Greet(nil)
	if err != nil {
		log.Fatalf("calling hello world contract '%s': %v", address, err)
	}

	log.Printf("Greeting from contract '%s': %s", address, greeting)
}
