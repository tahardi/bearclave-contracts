package foundry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	AnvilCommand  = "anvil"
	ForgeCommand  = "forge"
	ScriptCommand = "script"

	BroadcastFlag = "--broadcast"
	PrivateKeyFlag = "--private-key"
	RPCFlag        = "--rpc-url"

	// BroadcastPath is where the `forge script` stores the resulting broadcast file.
	// It should be: <broadcast_dir>/<script_name>/<chain_id>/run-latest.json
	//
	// Example: ../../contracts/broadcast/HelloWorld.s.sol/31337/run-latest.json
	BroadcastPath = "%s/%s/%d/run-latest.json"

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

	BaseFee          = 1_000_000_000
	ChainID          = 31337
	GasLimit         = 30_000_000
	GenesisTimestamp = 1769011998
	GenesisNumber    = 0
	StartingBalance  = 10_000
	URL              = "http://127.0.0.1:8545"
	WaitTime         = 250 * time.Millisecond
)

var (
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
	accounts, err := NewDefaultAnvilAccounts()
	if err != nil {
		return nil, fmt.Errorf("%w: creating accounts: %w", ErrAnvil, err)
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

func (a *Anvil) Start(ctx context.Context) error {
	go func() {
		a.server = exec.CommandContext(ctx, AnvilCommand)
		out, err := a.server.CombinedOutput()
		if err != nil {
			fmt.Printf("%s: starting anvil: %s\n", ErrAnvil.Error(), string(out))
		}
	}()

	time.Sleep(WaitTime)
	return nil
}

func (a *Anvil) Stop() error {
	if a.server != nil {
		return a.server.Process.Kill()
	}
	return nil
}

// DeployContract deploys a smart contract via the `forge script` command.
// The command format is:
// forge script <script_path> --rpc-url <rpc_url> --private-key <private_key> --broadcast
//
// Example:
//
//	forge script ../../contracts/scripts/HelloWorld.s.sol:HelloWorldScript \
//	    --rpc-url http://127.0.0.1:8545 \
//	    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
//	    --broadcast
func (a *Anvil) DeployContract(
	ctx context.Context,
	contractName string,
	owner *Account,
) (*common.Address, error) {
	scriptName := fmt.Sprintf(ScriptName, contractName)
	scriptPath := fmt.Sprintf(ScriptPath, a.scriptDir, scriptName, contractName)
	args := []string{
		ScriptCommand,
		scriptPath,
		RPCFlag, a.url,
		PrivateKeyFlag, owner.PrivateKeyHex(),
		BroadcastFlag,
	}

	deploy := exec.CommandContext(ctx, ForgeCommand, args...)
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
