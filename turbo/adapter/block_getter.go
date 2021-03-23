package adapter

import (
	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

func NewBlockGetter(dbReader ethdb.Database) *blockGetter {
	return &blockGetter{dbReader}
}

type blockGetter struct {
	dbReader ethdb.Database
}

func (g *blockGetter) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	return rawdb.ReadBlockByHash(g.dbReader, hash)
}

func (g *blockGetter) GetBlock(hash common.Hash, number uint64) *types.Block {
	return rawdb.ReadBlock(g.dbReader, hash, number)
}
