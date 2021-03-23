package commands

import (
	"context"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/forkid"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

// Forks is a data type to record a list of forks passed by this node
type Forks struct {
	GenesisHash common.Hash `json:"genesis"`
	Forks       []uint64    `json:"forks"`
}

// Forks implements tg_forks. Returns the genesis block hash and a sorted list of all forks block numbers
func (api *TgImpl) Forks(ctx context.Context) (Forks, error) {
	tx, err := api.db.Begin(ctx, ethdb.RO)
	if err != nil {
		return Forks{}, err
	}
	defer tx.Rollback()

	chainConfig, genesis, err := api.chainConfigWithGenesis(tx)
	if err != nil {
		return Forks{}, err
	}
	forksBlocks := forkid.GatherForks(chainConfig)

	return Forks{genesis.Hash(), forksBlocks}, nil
}
