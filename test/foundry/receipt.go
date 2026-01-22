package foundry

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrReceipt = errors.New("receipt")
)

type Receipt struct {
	Status            uint64
	CumulativeGasUsed uint64
	Logs              []*Log
	LogsBloom         []byte
	Type              uint64
	TransactionHash   []byte
	TransactionIndex  uint64
	BlockHash         []byte
	BlockNumber       uint64
	GasUsed           uint64
	EffectiveGasPrice uint64
	BlobGasPrice      uint64
	From              *common.Address
	To                *common.Address
	ContractAddress   *common.Address
}

//nolint:tagliatelle
type receiptJSON struct {
	Status            string `json:"status"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Logs              []*Log `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Type              string `json:"type"`
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	GasUsed           string `json:"gasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	BlobGasPrice      string `json:"blobGasPrice"`
	From              *common.Address `json:"from"`
	To                *common.Address `json:"to"`
	ContractAddress   *common.Address `json:"contractAddress"`
}

func (r *Receipt) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiptJSON{
		Status:            Uint64ToHexString(r.Status),
		CumulativeGasUsed: Uint64ToHexString(r.CumulativeGasUsed),
		Logs:              r.Logs,
		LogsBloom:         BytesToHexString(r.LogsBloom),
		Type:              Uint64ToHexString(r.Type),
		TransactionHash:   BytesToHexString(r.TransactionHash),
		TransactionIndex:  Uint64ToHexString(r.TransactionIndex),
		BlockHash:         BytesToHexString(r.BlockHash),
		BlockNumber:       Uint64ToHexString(r.BlockNumber),
		GasUsed:           Uint64ToHexString(r.GasUsed),
		EffectiveGasPrice: Uint64ToHexString(r.EffectiveGasPrice),
		BlobGasPrice:      Uint64ToHexString(r.BlobGasPrice),
		From:              r.From,
		To:                r.To,
		ContractAddress:   r.ContractAddress,
	})
}

func (r *Receipt) UnmarshalJSON(data []byte) error {
	receipt := receiptJSON{}
	err := json.Unmarshal(data, &receipt)
	if err != nil {
		return fmt.Errorf("%w: unmarshaling receipt: %w", ErrReceipt, err)
	}

	status, err := ParseUint64FromHexString(receipt.Status)
	if err != nil {
		return fmt.Errorf("%w: parsing status: %w", ErrReceipt, err)
	}

	cumulativeGasUsed, err := ParseUint64FromHexString(receipt.CumulativeGasUsed)
	if err != nil {
		return fmt.Errorf("%w: parsing cumulative gas used: %w", ErrReceipt, err)
	}

	logsBloom, err := ParseBytesFromHexString(receipt.LogsBloom)
	if err != nil {
		return fmt.Errorf("%w: decoding logs bloom: %w", ErrReceipt, err)
	}

	rType, err := ParseUint64FromHexString(receipt.Type)
	if err != nil {
		return fmt.Errorf("%w: parsing receipt type: %w", ErrReceipt, err)
	}

	transactionHash, err := ParseBytesFromHexString(receipt.TransactionHash)
	if err != nil {
		return fmt.Errorf("%w: decoding transaction hash: %w", ErrReceipt, err)
	}

	transactionIndex, err := ParseUint64FromHexString(receipt.TransactionIndex)
	if err != nil {
		return fmt.Errorf("%w: parsing transaction index: %w", ErrReceipt, err)
	}

	blockHash, err := ParseBytesFromHexString(receipt.BlockHash)
	if err != nil {
		return fmt.Errorf("%w: parsing block hash: %w", ErrReceipt, err)
	}

	blockNumber, err := ParseUint64FromHexString(receipt.BlockNumber)
	if err != nil {
		return fmt.Errorf("%w: parsing block number: %w", ErrReceipt, err)
	}

	gasUsed, err := ParseUint64FromHexString(receipt.GasUsed)
	if err != nil {
		return fmt.Errorf("%w: parsing gas used: %w", ErrReceipt, err)
	}

	effectiveGasPrice, err := ParseUint64FromHexString(receipt.EffectiveGasPrice)
	if err != nil {
		return fmt.Errorf("%w: parsing effective gas price: %w", ErrReceipt, err)
	}

	blobGasPrice, err := ParseUint64FromHexString(receipt.BlobGasPrice)
	if err != nil {
		return fmt.Errorf("%w: parsing blob gas price: %w", ErrReceipt, err)
	}

	r.Status = status
	r.CumulativeGasUsed = cumulativeGasUsed
	r.Logs = receipt.Logs
	r.LogsBloom = logsBloom
	r.Type = rType
	r.TransactionHash = transactionHash
	r.TransactionIndex = transactionIndex
	r.BlockHash = blockHash
	r.BlockNumber = blockNumber
	r.GasUsed = gasUsed
	r.EffectiveGasPrice = effectiveGasPrice
	r.BlobGasPrice = blobGasPrice
	r.From = receipt.From
	r.To = receipt.To
	r.ContractAddress = receipt.ContractAddress
	return nil
}
