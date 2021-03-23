package commands

import (
	"context"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/rpc"
)

// TgAPI TurboGeth specific routines
type TgAPI interface {
	// System related (see ./tg_system.go)
	Forks(ctx context.Context) (Forks, error)

	// Blocks related (see ./tg_blocks.go)
	GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	GetHeaderByHash(_ context.Context, hash common.Hash) (*types.Header, error)

	// Receipt related (see ./tg_receipts.go)
	GetLogsByHash(ctx context.Context, hash common.Hash) ([][]*types.Log, error)
	//GetLogsByNumber(ctx context.Context, number rpc.BlockNumber) ([][]*types.Log, error)

	// Issuance / reward related (see ./tg_issuance.go)
	// BlockReward(ctx context.Context, blockNr rpc.BlockNumber) (Issuance, error)
	// UncleReward(ctx context.Context, blockNr rpc.BlockNumber) (Issuance, error)
	Issuance(ctx context.Context, blockNr rpc.BlockNumber) (Issuance, error)
}

// TgImpl is implementation of the TgAPI interface
type TgImpl struct {
	*BaseAPI
	db ethdb.Database
}

// NewTgAPI returns TgImpl instance
func NewTgAPI(db ethdb.Database) *TgImpl {
	return &TgImpl{
		BaseAPI: &BaseAPI{},
		db:      db,
	}
}
