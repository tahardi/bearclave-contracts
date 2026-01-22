package foundry

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	StartCommand = "anvil"

	// ScriptName is the name of the file containing the scrip to deploy a contract.
	// It should be: <contract_name>.s.sol
	//
	// Example: HelloWorld.s.sol
	ScriptName = "%s.s.sol"

	// ScriptPath is the script path to use with the `forge script` command.
	// It should be: <script_dir>/<script_name>:<contract_name>Script
	//
	// Example: ../../contracts/scripts/HelloWorld.s.sol:HelloWorldScript
	ScriptPath = "%s/%s:%sScript"

	// DeployCommand deploys a smart contract via the `forge script` command.
	// It should be: <script_path> --rpc-url <rpc_url> --private-key <private_key> --broadcast
	//
	// Example:
	// forge script ../../contracts/scripts/HelloWorld.s.sol:HelloWorldScript \
	//     --rpc-url http://127.0.0.1:8545 \
	//     --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
	//     --broadcast
	DeployCommand = "forge script %s --rpc-url %s --private-key %s --broadcast"

	// BroadcastPath is where the `forge script` stores the result broadcast file.
	// It should be: <broadcast_dir>/<script_name>/<chain_id>/run-latest.json
	//
	// Example: ../../contracts/broadcast/HelloWorld.s.sol/31337/run-latest.json
	BroadcastPath = "%s/%s/%d/run-latest.json"

	Address0  = "0x0000000000000000000000000000000000000000"
	Address1  = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	Address2  = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	Address3  = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
	Address4  = "0x90F79bf6EB2c4f870365E785982E1f101E93b906"
	Address5  = "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
	Address6  = "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"
	Address7  = "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
	Address8  = "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"
	Address9  = "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f"
	Address10 = "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720"

	PrivateKey1  = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	PrivateKey2  = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	PrivateKey3  = "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
	PrivateKey4  = "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"
	PrivateKey5  = "0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a"
	PrivateKey6  = "0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba"
	PrivateKey7  = "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e"
	PrivateKey8  = "0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356"
	PrivateKey9  = "0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97"
	PrivateKey10 = "0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6"

	Mnemonic       = "test test test test test test test test test test test junk"
	DerivationPath = "m/44'/60'/0'/0/0"

	BaseFee          = 1_000_000_000
	ChainID          = 31337
	GasLimit         = 30_000_000
	GenesisTimestamp = 1769011998
	GenesisNumber    = 0
	StartingBalance  = 10_000

	URL = "http://127.0.0.1:8545"
)

var (
	Addresses = []string{
		Address1, Address2, Address3, Address4, Address5,
		Address6, Address7, Address8, Address9, Address10,
	}
	PrivateKeys = []string{
		PrivateKey1, PrivateKey2, PrivateKey3, PrivateKey4, PrivateKey5,
		PrivateKey6, PrivateKey7, PrivateKey8, PrivateKey9, PrivateKey10,
	}
	ErrAnvil = errors.New("anvil")
)

type Anvil struct {
	accounts         []*Account
	baseFee          uint64
	chainID          uint64
	gasLimit         uint64
	genesisTimestamp uint64
	genesisNumber    uint64
	url              string
	server           *exec.Cmd
	broadcastDir     string
	scriptDir        string
}

func NewAnvil(
	broadcastDir string,
	scriptDir string,
) (*Anvil, error) {
	accounts := make([]*Account, len(Addresses))
	for i, address := range Addresses {
		account, err := NewAccount(address, PrivateKeys[i], StartingBalance)
		if err != nil {
			return nil, fmt.Errorf("%w: making account: %w"+
				"", ErrAnvil, err)
		}
		accounts[i] = account
	}

	return &Anvil{
		accounts:         accounts,
		baseFee:          BaseFee,
		chainID:          ChainID,
		gasLimit:         GasLimit,
		genesisTimestamp: GenesisTimestamp,
		genesisNumber:    GenesisNumber,
		url:              URL,
		broadcastDir:     broadcastDir,
		scriptDir:        scriptDir,
	}, nil
}

func (a *Anvil) Accounts() []*Account     { return a.accounts }
func (a *Anvil) Account(i int) *Account   { return a.accounts[i] }
func (a *Anvil) BaseFee() uint64          { return a.baseFee }
func (a *Anvil) ChainID() uint64          { return a.chainID }
func (a *Anvil) GasLimit() uint64         { return a.gasLimit }
func (a *Anvil) GenesisTimestamp() uint64 { return a.genesisTimestamp }
func (a *Anvil) GenesisNumber() uint64    { return a.genesisNumber }
func (a *Anvil) URL() string              { return a.url }

func (a *Anvil) Start() error {
	go func() {
		a.server = exec.Command(StartCommand)
		out, err := a.server.CombinedOutput()
		if err != nil {
			fmt.Printf("%w: starting anvil: %s\n", ErrAnvil, string(out))
		}
	}()

	time.Sleep(time.Millisecond * 250)
	return nil
}

func (a *Anvil) Stop() error {
	if a.server != nil {
		return a.server.Process.Kill()
	}
	return nil
}

func (a *Anvil) DeployContract(
	contractName string,
	owner *Account,
) (*common.Address, error) {
	scriptName := fmt.Sprintf(ScriptName, contractName)
	scriptPath := fmt.Sprintf(ScriptPath, a.scriptDir, scriptName, contractName)
	command := fmt.Sprintf(DeployCommand, scriptPath, a.URL(), owner.PrivateKeyHex())
	fields := strings.Fields(command)

	deploy := exec.Command(fields[0], fields[1:]...)
	out, err := deploy.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%w: deploying contract: %w: %s", ErrAnvil, err, string(out))
	}

	broadcastPath := fmt.Sprintf(BroadcastPath, a.broadcastDir, scriptName, a.chainID)
	bytes, err := os.ReadFile(broadcastPath)
	if err != nil {
		return nil, fmt.Errorf("%w: reading broadcast file: %w", ErrAnvil, err)
	}

	broadcast := &Broadcast{}
	err = json.Unmarshal(bytes, broadcast)
	if err != nil {
		return nil, fmt.Errorf("%w: unmarshaling broadcast: %w", ErrAnvil, err)
	}

	address, err := broadcast.GetContractAddress(contractName)
	if err != nil {
		return nil, fmt.Errorf("%w: getting contract address: %w", ErrAnvil, err)
	}
	return address, nil
}
