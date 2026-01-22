package foundry

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	DefaultAddress1 = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	DefaultAddress2 = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	DefaultAddress3 = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
	DefaultAddress4 = "0x90F79bf6EB2c4f870365E785982E1f101E93b906"
	DefaultAddress5 = "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
	DefaultAddress6 = "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"
	DefaultAddress7 = "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
	DefaultAddress8 = "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"
	DefaultAddress9  = "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f"
	DefaultAddress10 = "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720"

	DefaultPrivateKey1 = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	DefaultPrivateKey2 = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	DefaultPrivateKey3 = "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
	DefaultPrivateKey4 = "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"
	DefaultPrivateKey5 = "0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a"
	DefaultPrivateKey6 = "0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba"
	DefaultPrivateKey7 = "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e"
	DefaultPrivateKey8 = "0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356"
	DefaultPrivateKey9  = "0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97"
	DefaultPrivateKey10 = "0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6"
)

var (
	DefaultAddresses = []string{
		DefaultAddress1, DefaultAddress2, DefaultAddress3, DefaultAddress4, DefaultAddress5,
		DefaultAddress6, DefaultAddress7, DefaultAddress8, DefaultAddress9, DefaultAddress10,
	}
	DefaultPrivateKeys = []string{
		DefaultPrivateKey1, DefaultPrivateKey2, DefaultPrivateKey3, DefaultPrivateKey4, DefaultPrivateKey5,
		DefaultPrivateKey6, DefaultPrivateKey7, DefaultPrivateKey8, DefaultPrivateKey9, DefaultPrivateKey10,
	}

	ErrAccount = errors.New("account")
)

type Account struct {
	address       *common.Address
	privateKey    *ecdsa.PrivateKey
	privateKeyHex string
	balance       uint64
}

func NewDefaultAnvilAccounts() ([]*Account, error) {
	accounts := make([]*Account, len(DefaultAddresses))
	for i, address := range DefaultAddresses {
		account, err := NewAccount(address, DefaultPrivateKeys[i], StartingBalance)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to create account: %w", ErrAccount, err)
		}
		accounts[i] = account
	}
	return accounts, nil
}

func NewAccount(
	address string,
	privateKeyHex string,
	balance uint64,
) (*Account, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, HexStringPrefix))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid private key: %w", ErrAccount, err)
	}

	addr := common.HexToAddress(address)
	return &Account{
		address:       &addr,
		privateKey:    privateKey,
		privateKeyHex: privateKeyHex,
		balance:       balance,
	}, nil
}

func (a *Account) Address() *common.Address      { return a.address }
func (a *Account) PrivateKey() *ecdsa.PrivateKey { return a.privateKey }
func (a *Account) PrivateKeyHex() string         { return a.privateKeyHex }
