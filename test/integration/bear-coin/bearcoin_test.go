package bearcoin_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/test/integration"
)

const (
	ContractName = "BearCoin"
	Decimals     = 18
	Base         = 1_000_000
)

// func TestBearCoin_Allowance(t *testing.T) {
//
// }
//
// func TestBearCoin_Approve(t *testing.T) {
//
// }

func TestBearCoin_BalanceOf(t *testing.T) {
	t.Run("happy path - owner", func(t *testing.T) {
		// given
		want := totalSupply()

		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

		// when
		got, err := contract.BalanceOf(nil, owner.Address())

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - other", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

		other := anvil.Account(1)

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
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		burnAmount := big.NewInt(100)
		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)
		requireBalance(t, contract, owner, totalSupply())

		// when
		_, err := burn(t, anvil, contract, owner, burnAmount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), burnAmount))
	})

	t.Run("happy path - other burn", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := anvil.Account(0), anvil.Account(1)
		contract := deployContract(t, anvil, owner)

		_, err := burn(t, anvil, contract, owner, amount)
		require.NoError(t, err)

		_, err = mint(t, anvil, contract, owner, other, amount)
		require.NoError(t, err)
		requireBalance(t, contract, other, amount)

		// when
		_, err = burn(t, anvil, contract, other, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, other, nil)
	})

	t.Run("error - burn amount greater than supply", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		burnAmount := totalSupply().Add(totalSupply(), big.NewInt(1))
		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)
		requireBalance(t, contract, owner, totalSupply())

		// when
		_, err := burn(t, anvil, contract, owner, burnAmount)

		// then
		require.Error(t, err)
	})

	t.Run("error - user has no tokens to burn", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

		burnAmount := totalSupply()
		brokeUser := anvil.Account(1)
		requireBalance(t, contract, brokeUser, nil)

		// when
		_, err := burn(t, anvil, contract, brokeUser, burnAmount)

		// then
		require.Error(t, err)
	})
}

func TestBearCoin_Mint(t *testing.T) {
	t.Run("happy path - owner mint-to-self", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)
		requireBalance(t, contract, owner, totalSupply())

		_, err := burn(t, anvil, contract, owner, amount)
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))

		// when
		_, err = mint(t, anvil, contract, owner, owner, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply())
	})

	t.Run("happy path - owner mint-to-other", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := anvil.Account(0), anvil.Account(1)
		contract := deployContract(t, anvil, owner)
		requireBalance(t, contract, owner, totalSupply())

		_, err := burn(t, anvil, contract, owner, amount)
		require.NoError(t, err)
		requireBalance(t, contract, owner, totalSupply().Sub(totalSupply(), amount))
		requireBalance(t, contract, other, nil)

		// when
		_, err = mint(t, anvil, contract, owner, other, amount)

		// then
		require.NoError(t, err)
		requireBalance(t, contract, other, amount)
	})

	t.Run("error - only owner can mint", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		amount := big.NewInt(100)
		owner, other := anvil.Account(0), anvil.Account(1)
		contract := deployContract(t, anvil, owner)

		// when
		_, err := mint(t, anvil, contract, other, other, amount)

		// then
		require.Error(t, err)
	})
}

func TestBearCoin_Name(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		want := "BearCoin"
		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

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
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		want := "BCN"
		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

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
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		want := totalSupply()
		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

		// when
		got, err := contract.TotalSupply(nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

// func TestBearCoin_Transfer(t *testing.T) {
//
// }
//
// func TestBearCoin_TransferFrom(t *testing.T) {
//
// }

func TestBearCoin_TransferOwnership(t *testing.T) {
	t.Run("happy path - deployer is owner", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		want := anvil.Account(0)
		contract := deployContract(t, anvil, want)

		// when
		got, err := contract.Owner(nil)

		// then
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, want.Address(), got)
	})

	t.Run("happy path - transfer ownership", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		oldOwner := anvil.Account(0)
		contract := deployContract(t, anvil, oldOwner)

		got, err := contract.Owner(nil)
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, oldOwner.Address(), got)

		newOwner := anvil.Account(1)
		opts := newTransactionOpts(t, anvil, oldOwner)
		call := func() (*types.Transaction, error) {
			return contract.TransferOwnership(opts, newOwner.Address())
		}

		// when
		_, err = executeCall(t, anvil, call)

		// then
		require.NoError(t, err)
		got, err = contract.Owner(nil)
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, newOwner.Address(), got)
	})

	t.Run("error - other cannot transfer Ownership", func(t *testing.T) {
		// given
		anvil, stop := integration.StartAnvil(t, true)
		defer stop()

		owner := anvil.Account(0)
		contract := deployContract(t, anvil, owner)

		got, err := contract.Owner(nil)
		require.NoError(t, err)
		integration.AssertAddressesEqual(t, owner.Address(), got)

		other := anvil.Account(1)
		opts := newTransactionOpts(t, anvil, other)
		call := func() (*types.Transaction, error) {
			return contract.TransferOwnership(opts, other.Address())
		}

		// when
		_, err = executeCall(t, anvil, call)

		// then
		require.Error(t, err)
	})
}
