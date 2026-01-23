package bearcoin_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/contracts/bindings"
	"github.com/tahardi/bearchain/test/foundry"
)

func approve(
	t *testing.T,
	anvil *foundry.Anvil,
	contract *bindings.BearCoin,
	principal *foundry.Account,
	proxy *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, anvil, principal)
	call := func() (*types.Transaction, error) {
		return contract.Approve(opts, proxy.Address(), amount)
	}
	return executeCall(t, anvil, call)
}

func burn(
	t *testing.T,
	anvil *foundry.Anvil,
	contract *bindings.BearCoin,
	account *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, anvil, account)
	call := func() (*types.Transaction, error) {
		return contract.Burn(opts, amount)
	}
	return executeCall(t, anvil, call)
}

func deployContract(
	t *testing.T,
	anvil *foundry.Anvil,
	owner *foundry.Account,
) *bindings.BearCoin {
	t.Helper()
	contractAddress, err := anvil.DeployContract(t.Context(), ContractName, owner)
	require.NoError(t, err)

	client, err := anvil.Client()
	require.NoError(t, err)

	contract, err := bindings.NewBearCoin(*contractAddress, client)
	require.NoError(t, err)
	return contract
}

func executeCall(
	t *testing.T,
	anvil *foundry.Anvil,
	contractCall func() (*types.Transaction, error),
) (*types.Receipt, error) {
	t.Helper()
	tx, err := contractCall()
	if err != nil {
		return nil, err
	}

	client, err := anvil.Client()
	if err != nil {
		return nil, err
	}
	return bind.WaitMined(t.Context(), client, tx)
}

func mint(
	t *testing.T,
	anvil *foundry.Anvil,
	contract *bindings.BearCoin,
	owner *foundry.Account,
	to *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, anvil, owner)
	call := func() (*types.Transaction, error) {
		return contract.Mint(opts, to.Address(), amount)
	}
	return executeCall(t, anvil, call)
}

func newTransactionOpts(
	t *testing.T,
	anvil *foundry.Anvil,
	from *foundry.Account,
) *bind.TransactOpts {
	t.Helper()
	opts, err := bind.NewKeyedTransactorWithChainID(from.PrivateKey(), anvil.ChainID())
	require.NoError(t, err)
	return opts
}

func requireAllowance(
	t *testing.T,
	contract *bindings.BearCoin,
	principal *foundry.Account,
	proxy *foundry.Account,
	want *big.Int,
) {
	t.Helper()
	got, err := contract.Allowance(nil, principal.Address(), proxy.Address())
	require.NoError(t, err)
	if want == nil {
		require.Equal(t, 0, got.Cmp(big.NewInt(0)))
	} else {
		require.Equal(t, want, got)
	}
}

func requireBalance(
	t *testing.T,
	contract *bindings.BearCoin,
	account *foundry.Account,
	want *big.Int,
) {
	t.Helper()
	got, err := contract.BalanceOf(nil, account.Address())
	require.NoError(t, err)
	if want == nil {
		require.Equal(t, 0, got.Cmp(big.NewInt(0)))
	} else {
		require.Equal(t, want, got)
	}
}

func requireMaxUint256(t *testing.T) *big.Int {
	t.Helper()
	maxUint256, ok := new(big.Int).
		SetString(
			"115792089237316195423570985008687907853269984665640564039457584007913129639935",
			10,
		)
	require.True(t, ok)
	return maxUint256
}

func totalSupply() *big.Int {
	decimals := big.NewInt(Decimals)
	base := big.NewInt(Base)
	ten := big.NewInt(10)
	pow := big.NewInt(0).Exp(ten, decimals, nil)
	return big.NewInt(0).Mul(pow, base)
}

func transfer(
	t *testing.T,
	anvil *foundry.Anvil,
	contract *bindings.BearCoin,
	from *foundry.Account,
	to *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, anvil, from)
	call := func() (*types.Transaction, error) {
		return contract.Transfer(opts, to.Address(), amount)
	}
	return executeCall(t, anvil, call)
}

func transferFrom(
	t *testing.T,
	anvil *foundry.Anvil,
	contract *bindings.BearCoin,
	principal *foundry.Account,
	proxy *foundry.Account,
	to *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, anvil, proxy)
	call := func() (*types.Transaction, error) {
		return contract.TransferFrom(opts, principal.Address(), to.Address(), amount)
	}
	return executeCall(t, anvil, call)
}
