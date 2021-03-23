package core

import (
	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb/remote"
	"github.com/EVRICE/tgeth_alpha/rlp"
)

type EthBackend struct {
	Backend
}

type Backend interface {
	TxPool() *TxPool
	Etherbase() (common.Address, error)
	NetVersion() (uint64, error)
}

func NewEthBackend(eth Backend) *EthBackend {
	return &EthBackend{eth}
}

func (back *EthBackend) AddLocal(signedtx []byte) ([]byte, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(signedtx, tx); err != nil {
		return common.Hash{}.Bytes(), err
	}

	return tx.Hash().Bytes(), back.TxPool().AddLocal(tx)
}

func (back *EthBackend) Subscribe(func(*remote.SubscribeReply)) error {
	// do nothing
	return nil
}
