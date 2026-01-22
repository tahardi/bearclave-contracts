package foundry

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrAccount = errors.New("account")
)

type Account struct {
	address       *common.Address
	privateKey    *ecdsa.PrivateKey
	privateKeyHex string
	balance       uint64
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
