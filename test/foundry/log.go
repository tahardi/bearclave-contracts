package foundry

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrLog = errors.New("log")
)

type Log struct {
	Address          string
	Topics           []string
	Data             string
	BlockHash        []byte
	BlockNumber      uint64
	BlockTimestamp   uint64
	TransactionHash  []byte
	TransactionIndex uint64
	LogIndex         uint64
	Removed          bool
}

//nolint:tagliatelle
type logJSON struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	BlockTimestamp   string   `json:"blockTimestamp"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func (l *Log) MarshalJSON() ([]byte, error) {
	return json.Marshal(logJSON{
		Address:          l.Address,
		Topics:           l.Topics,
		Data:             l.Data,
		BlockHash:        BytesToHexString(l.BlockHash),
		BlockNumber:      Uint64ToHexString(l.BlockNumber),
		BlockTimestamp:   Uint64ToHexString(l.BlockTimestamp),
		TransactionHash:  BytesToHexString(l.TransactionHash),
		TransactionIndex: Uint64ToHexString(l.TransactionIndex),
		LogIndex:         Uint64ToHexString(l.LogIndex),
		Removed:          l.Removed,
	})
}

func (l *Log) UnmarshalJSON(data []byte) error {
	log := logJSON{}
	err := json.Unmarshal(data, &log)
	if err != nil {
		return fmt.Errorf("%w: unmarshaling log: %w", ErrLog, err)
	}

	blockHash, err := ParseBytesFromHexString(log.BlockHash)
	if err != nil {
		return fmt.Errorf("%w: decoding block hash: %w", ErrLog, err)
	}

	blockNumber, err := ParseUint64FromHexString(log.BlockNumber)
	if err != nil {
		return fmt.Errorf("%w: parsing block number: %w", ErrLog, err)
	}

	blockTimestamp, err := ParseUint64FromHexString(log.BlockTimestamp)
	if err != nil {
		return fmt.Errorf("%w: parsing block timestamp: %w", ErrLog, err)
	}

	transactionHash, err := ParseBytesFromHexString(log.TransactionHash)
	if err != nil {
		return fmt.Errorf("%w: decoding transaction hash: %w", ErrLog, err)
	}

	transactionIndex, err := ParseUint64FromHexString(log.TransactionIndex)
	if err != nil {
		return fmt.Errorf("%w: parsing transaction index: %w", ErrLog, err)
	}

	logIndex, err := ParseUint64FromHexString(log.LogIndex)
	if err != nil {
		return fmt.Errorf("%w: parsing log index: %w", ErrLog, err)
	}

	l.Address = log.Address
	l.Topics = log.Topics
	l.Data = log.Data
	l.BlockHash = blockHash
	l.BlockNumber = blockNumber
	l.BlockTimestamp = blockTimestamp
	l.TransactionHash = transactionHash
	l.TransactionIndex = transactionIndex
	l.LogIndex = logIndex
	l.Removed = log.Removed
	return nil
}
