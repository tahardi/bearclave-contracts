package isevenverifier_test

import (
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-contracts/contracts/bindings"
	"github.com/tahardi/bearclave-foundry/foundry"
	"github.com/tahardi/bearclave-contracts/test/integration"
)

const (
	ContractName = "FactorsVerifier"
	// This is the hex representation of the base64 image_id string from
	// the testdata/groth16.json file. This hex string is used in the
	// contracts/scripts/FactorsVerifier.s.sol file for deploying the contract.
	// This is here to detect updates to groth16.json that break the script.
	ImageID = "2bda51ae4f0326636e89480a383a770f713d514f76f6d58a5c2d1373b6b87d48"
)

//go:embed testdata/groth16.json
var groth16Bytes []byte

type Groth16Proof struct {
	ImageID []byte `json:"image_id"`
	Journal []byte `json:"journal"`
	Seal    []byte `json:"seal"`
}

func groth16ProofFromTestData(t *testing.T) Groth16Proof {
	t.Helper()
	groth16Proof := Groth16Proof{}
	err := json.Unmarshal(groth16Bytes, &groth16Proof)
	require.NoError(t, err)
	return groth16Proof
}

func imageIDFromHexString(t *testing.T, imageIDHex string) [32]byte {
	t.Helper()
	bytes, err := hex.DecodeString(imageIDHex)
	require.NoError(t, err)

	var imageID [32]byte
	copy(imageID[:], bytes)
	return imageID
}

func deployContract(
	t *testing.T,
	f *foundry.Foundry,
	owner *foundry.Account,
) *bindings.FactorsVerifier {
	t.Helper()
	contractAddress, err := f.Forge().DeployContract(t.Context(), ContractName, owner)
	require.NoError(t, err)

	client, err := f.Anvil().Client(t.Context())
	require.NoError(t, err)

	contract, err := bindings.NewFactorsVerifier(*contractAddress, client)
	require.NoError(t, err)
	return contract
}

func newTransactionOpts(
	t *testing.T,
	f *foundry.Foundry,
	from *foundry.Account,
) *bind.TransactOpts {
	t.Helper()
	opts, err := bind.NewKeyedTransactorWithChainID(from.PrivateKey(), f.Anvil().ChainID())
	require.NoError(t, err)
	return opts
}

func executeCall(
	t *testing.T,
	f *foundry.Foundry,
	contractCall func() (*types.Transaction, error),
) (*types.Receipt, error) {
	t.Helper()
	tx, err := contractCall()
	if err != nil {
		return nil, err
	}

	client, err := f.Anvil().Client(t.Context())
	if err != nil {
		return nil, err
	}
	return bind.WaitMined(t.Context(), client, tx)
}

func setImageID(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.FactorsVerifier,
	account *foundry.Account,
	imageID [32]byte,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, f, account)
	call := func() (*types.Transaction, error) {
		return contract.SetImageId(opts, imageID)
	}
	return executeCall(t, f, call)
}

func verify(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.FactorsVerifier,
	account *foundry.Account,
	product uint64,
	seal []byte,
) (*types.Receipt, error) {
	t.Helper()
	opts := newTransactionOpts(t, f, account)
	call := func() (*types.Transaction, error) {
		return contract.Verify(opts, product, seal)
	}
	return executeCall(t, f, call)
}

func requireFactorsKnownForProductEvent(
	t *testing.T,
	contract *bindings.FactorsVerifier,
	receipt *types.Receipt,
	wantProduct uint64,
) {
	t.Helper()
	for _, log := range receipt.Logs {
		event, err := contract.ParseFactorsKnownForProduct(*log)
		if err == nil {
			assert.Equal(t, wantProduct, event.Product)
			return
		}
	}
	t.Fatalf("expected FactorsKnownForProduct event but none was found")
}

func requireImageIDUpdatedEvent(
	t *testing.T,
	contract *bindings.FactorsVerifier,
	receipt *types.Receipt,
	wantImageID [32]byte,
) {
	t.Helper()
	for _, log := range receipt.Logs {
		event, err := contract.ParseImageIdUpdated(*log)
		if err == nil {
			assert.Equal(t, wantImageID, event.NewImageId)
			return
		}
	}
	t.Fatalf("expected ImageIdUpdated event but none was found")
}

func TestFactorsVerifier_SetImageID(t *testing.T) {
	t.Run("happy path - expected testdata image ID", func(t *testing.T) {
		// given
		want := imageIDFromHexString(t, ImageID)
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// when
		got, err := contract.ImageId(nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - change image ID", func(t *testing.T) {
		// given
		want := imageIDFromHexString(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// when
		receipt, err := setImageID(t, f, contract, owner, want)

		// then
		require.NoError(t, err)
		requireImageIDUpdatedEvent(t, contract, receipt, want)

		got, err := contract.ImageId(nil)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("error - only owner can set image ID", func(t *testing.T) {
		// given
		want := imageIDFromHexString(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner, other := f.Anvil().Account(0), f.Anvil().Account(1)
		contract := deployContract(t, f, owner)

		// when
		_, err := setImageID(t, f, contract, other, want)

		// then
		require.Error(t, err)
	})
}

func TestFactorsVerifier_Verify(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// Our Rust ZKVM program outputs the journal (i.e., our expected product)
		// in Big Endian format. If you want the actual value, convert to LE.
		groth16Proof := groth16ProofFromTestData(t)
		product := binary.BigEndian.Uint64(groth16Proof.Journal)

		// when
		receipt, err := verify(t, f, contract, owner, product, groth16Proof.Seal)

		// then
		require.NoError(t, err)
		requireFactorsKnownForProductEvent(t, contract, receipt, product)
	})

	t.Run("error - incorrect product", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		groth16Proof := groth16ProofFromTestData(t)
		product := uint64(8)

		// when
		_, err := verify(t, f, contract, owner, product, groth16Proof.Seal)

		// then
		require.Error(t, err)
	})

	t.Run("error - invalid seal", func(t *testing.T) {
		// given
		f, stop := integration.StartFoundry(t, true)
		defer stop()

		owner := f.Anvil().Account(0)
		contract := deployContract(t, f, owner)

		// Our Rust ZKVM program outputs the journal (i.e., our expected product)
		// in Big Endian format. If you want the actual value, convert to LE.
		groth16Proof := groth16ProofFromTestData(t)
		product := binary.BigEndian.Uint64(groth16Proof.Journal)
		groth16Proof.Seal = []byte("invalid seal")

		// when
		_, err := verify(t, f, contract, owner, product, groth16Proof.Seal)

		// then
		require.Error(t, err)
	})
}
