package bearcoin_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-contracts/test/integration"
)

const (
	ContractName = "BearCoin"
	Decimals     = 18
	Base         = 1_000_000
)

func TestBearCoin_Allowance(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)

		// when/then
		requireAllowance(t, contract, owner, other, nil)
		requireAllowance(t, contract, other, owner, nil)
	})
}

func TestBearCoin_Approve(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)
		requireAllowance(t, contract, owner, other, nil)

		// when
		_, err := approve(t, f, contract, owner, other, amount)

		// then
		require.NoError(t, err)
		requireAllowance(t, contract, owner, other, amount)
	})
}

func TestBearCoin_BalanceOf(t *testing.T) {
	t.Run("happy path - owner", func(t *testing.T) {
		// given
		want := totalSupply()

		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// when
		got, err := contract.BalanceOf(nil, owner.Address())

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - other", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		other := f.Anvil().Account(1)

		// when
		got, err := contract.BalanceOf(nil, other.Address())

		// then
		require.NoError(t, err)
		assert.Equal(t, 0, got.Cmp(big.NewInt(0)))
	})
}

func TestBearCoin_Burn(t *testing.T) {
	t.Run("happy path - owner burn", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		burnAmount := big.NewInt(100)
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())

		// when
		receipt, err := burn(t, f, contract, owner, burnAmount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), burnAmount))
		requireBurnEvent(t, contract, receipt, owner.Address(), burnAmount)
	})

	t.Run("happy path - other burn", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)

		_, err := burn(t, f, contract, owner, amount)
		require.NoError(t, err)

		_, err = mint(t, f, contract, owner, other, amount)
		require.NoError(t, err)
		requireBalance(t, contract, other, amount)

		// when
		receipt, err := burn(t, f, contract, other, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, other, nil)
		requireBurnEvent(t, contract, receipt, other.Address(), amount)
	})

	t.Run("error - burn amount greater than supply", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		burnAmount := totalSupply().Add(totalSupply(), big.NewInt(1))
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())

		// when
		_, err := burn(t, f, contract, owner, burnAmount)

		// then
		require.Error(t, err)
	})

	t.Run("error - user has no tokens to burn", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		burnAmount := totalSupply()
		brokeUser := f.Anvil().Account(1)
		requireBalance(t, contract, brokeUser, nil)

		// when
		_, err := burn(t, f, contract, brokeUser, burnAmount)

		// then
		require.Error(t, err)
	})
}

func TestBearCoin_Mint(t *testing.T) {
	t.Run("happy path - owner mint-to-self", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())

		_, err := burn(t, f, contract, owner, amount)
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))

		// when
		receipt, err := mint(t, f, contract, owner, owner, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply())
		requireMintEvent(t, contract, receipt, owner.Address(), amount)
	})

	t.Run("happy path - owner mint-to-other", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())

		_, err := burn(t, f, contract, owner, amount)
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))
		requireBalance(t, contract, other, nil)

		// when
		receipt, err := mint(t, f, contract, owner, other, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, other, amount)
		requireMintEvent(t, contract, receipt, other.Address(), amount)
	})

	t.Run("error - only owner can mint", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)

		// when
		_, err := mint(t, f, contract, other, other, amount)

		// then
		require.Error(t, err)
	})
}

func TestBearCoin_Name(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		want := "BearCoin"
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// when
		got, err := contract.Name(nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestBearCoin_Symbol(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		want := "BCN"
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// when
		got, err := contract.Symbol(nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestBearCoin_TotalSupply(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		want := totalSupply()
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// when
		got, err := contract.TotalSupply(nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestBearCoin_Transfer(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, other, nil)

		// when
		_, err := transfer(t, f, contract, owner, other, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))
		requireBalance(t, contract, other, amount)
	})

	t.Run("happy path - transfer to self", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())

		// when
		_, err := transfer(t, f, contract, owner, owner, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply())
	})

	t.Run("error - insufficient funds", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, other, nil)

		// when
		_, err := transfer(t, f, contract, other, owner, amount)

		// then
		require.Error(t, err)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, other, nil)
	})
}

func TestBearCoin_TransferFrom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, alice, bob := f.Anvil().Account(0), f.Anvil().Account(1), f.Anvil().Account(2)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, bob, nil)

		_, err := approve(t, f, contract, owner, alice, amount)
		require.NoError(t, err)
		requireAllowance(t, contract, owner, alice, amount)

		// when
		_, err = transferFrom(t, f, contract, owner, alice, bob, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))
		requireBalance(t, contract, bob, amount)
		requireAllowance(t, contract, owner, alice, nil)
	})

	t.Run("happy path - unlimited allowance", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		unlimited := requireMaxUint256(t)
		amount := big.NewInt(100)
		owner, alice, bob := f.Anvil().Account(0), f.Anvil().Account(1), f.Anvil().Account(2)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, bob, nil)

		_, err := approve(t, f, contract, owner, alice, unlimited)
		require.NoError(t, err)
		requireAllowance(t, contract, owner, alice, unlimited)

		// when
		_, err = transferFrom(t, f, contract, owner, alice, bob, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))
		requireBalance(t, contract, bob, amount)
		requireAllowance(t, contract, owner, alice, unlimited)
	})

	t.Run("error - insufficient allowance", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, alice, bob := f.Anvil().Account(0), f.Anvil().Account(1), f.Anvil().Account(2)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, bob, nil)
		requireAllowance(t, contract, owner, alice, nil)

		// when
		_, err := transferFrom(t, f, contract, owner, alice, bob, amount)

		// then
		require.Error(t, err)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, bob, nil)
		requireAllowance(t, contract, owner, alice, nil)
	})

	t.Run("error - insufficient funds", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		insufficient, amount := big.NewInt(100), big.NewInt(200)
		owner, alice, bob := f.Anvil().Account(0), f.Anvil().Account(1), f.Anvil().Account(2)
		contract := deployContract(t, f, owner)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, bob, nil)

		_, err := approve(t, f, contract, owner, alice, insufficient)
		require.NoError(t, err)
		requireAllowance(t, contract, owner, alice, insufficient)

		// when
		_, err = transferFrom(t, f, contract, owner, alice, bob, amount)

		// then
		require.Error(t, err)
		requireBalance(t, contract, owner, totalSupply())
		requireBalance(t, contract, bob, nil)
		requireAllowance(t, contract, owner, alice, insufficient)
	})
}

func TestBearCoin_TransferOwnership(t *testing.T) {
	t.Run("happy path - deployer is owner", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		want := f.Anvil().Account(0)
		contract := deployContract(t, f, want)

		// when
		got, err := contract.Owner(nil)

		// then
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, want.Address(), got)
	})

	t.Run("happy path - transfer ownership", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		oldOwner := f.Anvil().Account(0)
		contract := deployContract(t, f, oldOwner)

		got, err := contract.Owner(nil)
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, oldOwner.Address(), got)

		newOwner := f.Anvil().Account(1)
		opts := newTransactionOpts(t, f, oldOwner)
		call := func() (*types.Transaction, error) {
			return contract.TransferOwnership(opts, newOwner.Address())
		}

		// when
		_, err = executeCall(t, f, call)

		// then
		require.NoError(t, err)
		got, err = contract.Owner(nil)
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, newOwner.Address(), got)
	})

	t.Run("error - other cannot transfer Ownership", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		got, err := contract.Owner(nil)
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, owner.Address(), got)

		other := f.Anvil().Account(1)
		opts := newTransactionOpts(t, f, other)
		call := func() (*types.Transaction, error) {
			return contract.TransferOwnership(opts, other.Address())
		}

		// when
		_, err = executeCall(t, f, call)

		// then
		require.Error(t, err)
	})
}
