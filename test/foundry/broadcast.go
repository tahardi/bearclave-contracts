package foundry

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrBroadcast = errors.New("broadcast")
	ErrContractNotFound = fmt.Errorf("%w: contract not found", ErrBroadcast)
)

type Broadcast struct {
	Transactions []*Transaction `json:"transactions"`
	Receipts     []*Receipt     `json:"receipts"`
	Timestamp    uint64         `json:"timestamp"`
	Chain        uint64         `json:"chain"`
	Commit       string         `json:"commit"`
}

func (b *Broadcast) GetContractAddress(name string) (*common.Address, error) {
	for _, tx := range b.Transactions {
		if tx.ContractName == name {
			return tx.ContractAddress, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrContractNotFound, name)
}
