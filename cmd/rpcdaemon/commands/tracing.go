package commands

import (
	"context"
	"fmt"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/core/state"
	"github.com/EVRICE/tgeth_alpha/eth/tracers"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/internal/ethapi"
	"github.com/EVRICE/tgeth_alpha/rpc"
	"github.com/EVRICE/tgeth_alpha/turbo/adapter"
	"github.com/EVRICE/tgeth_alpha/turbo/rpchelper"
	"github.com/EVRICE/tgeth_alpha/turbo/transactions"
)

// TraceTransaction implements debug_traceTransaction. Returns Geth style transaction traces.
func (api *PrivateDebugAPIImpl) TraceTransaction(ctx context.Context, hash common.Hash, config *tracers.TraceConfig) (interface{}, error) {
	tx, err := api.dbReader.Begin(ctx, ethdb.RO)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Retrieve the transaction and assemble its EVM context
	txn, blockHash, _, txIndex := rawdb.ReadTransaction(tx, hash)
	if txn == nil {
		return nil, fmt.Errorf("transaction %#x not found", hash)
	}
	getter := adapter.NewBlockGetter(tx)
	chainContext := adapter.NewChainContext(tx)

	chainConfig, err := api.chainConfig(tx)
	if err != nil {
		return nil, err
	}

	msg, blockCtx, txCtx, ibs, _, err := transactions.ComputeTxEnv(ctx, getter, chainConfig, chainContext, tx.(ethdb.HasTx).Tx(), blockHash, txIndex)
	if err != nil {
		return nil, err
	}
	// Trace the transaction and return
	return transactions.TraceTx(ctx, msg, blockCtx, txCtx, ibs, config, chainConfig)
}

func (api *PrivateDebugAPIImpl) TraceCall(ctx context.Context, args ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *tracers.TraceConfig) (interface{}, error) {
	dbtx, err := api.dbReader.Begin(ctx, ethdb.RO)
	if err != nil {
		return nil, err
	}
	defer dbtx.Rollback()

	chainConfig, err := api.chainConfig(dbtx)
	if err != nil {
		return nil, err
	}

	blockNumber, hash, err := rpchelper.GetBlockNumber(blockNrOrHash, dbtx)
	if err != nil {
		return nil, err
	}
	var stateReader state.StateReader
	if num, ok := blockNrOrHash.Number(); ok && num == rpc.LatestBlockNumber {
		stateReader = state.NewPlainStateReader(dbtx)
	} else {
		stateReader = state.NewPlainDBState(dbtx, blockNumber)
	}
	header := rawdb.ReadHeader(dbtx, hash, blockNumber)
	if header == nil {
		return nil, fmt.Errorf("block %d(%x) not found", blockNumber, hash)
	}
	ibs := state.New(stateReader)
	msg := args.ToMessage(api.GasCap)
	blockCtx, txCtx := transactions.GetEvmContext(msg, header, blockNrOrHash.RequireCanonical, dbtx)
	// Trace the transaction and return
	return transactions.TraceTx(ctx, msg, blockCtx, txCtx, ibs, config, chainConfig)
}
