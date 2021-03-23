package commands

import (
	"context"
	"fmt"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/rpc"
)

// GetHeaderByNumber implements tg_getHeaderByNumber. Returns a block's header given a block number ignoring the block's transaction and uncle list (may be faster).
func (api *TgImpl) GetHeaderByNumber(ctx context.Context, blockNumber rpc.BlockNumber) (*types.Header, error) {
	tx, err := api.db.Begin(ctx, ethdb.RO)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	header := rawdb.ReadHeaderByNumber(tx, uint64(blockNumber.Int64()))
	if header == nil {
		return nil, fmt.Errorf("block header not found: %d", blockNumber.Int64())
	}

	return header, nil
}

// GetHeaderByHash implements tg_getHeaderByHash. Returns a block's header given a block's hash.
func (api *TgImpl) GetHeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	tx, err := api.db.Begin(ctx, ethdb.RO)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	header, err := rawdb.ReadHeaderByHash(tx, hash)
	if err != nil {
		return nil, err
	}
	if header == nil {
		return nil, fmt.Errorf("block header not found: %s", hash.String())
	}

	return header, nil
}
