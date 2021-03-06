package rpchelper

import (
	"fmt"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/core/types/accounts"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync/stages"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/rpc"
	"github.com/EVRICE/tgeth_alpha/turbo/adapter"
)

func GetBlockNumber(blockNrOrHash rpc.BlockNumberOrHash, dbReader ethdb.Database) (uint64, common.Hash, error) {
	var blockNumber uint64
	var err error
	hash, ok := blockNrOrHash.Hash()
	if !ok {
		number := *blockNrOrHash.BlockNumber
		if number == rpc.LatestBlockNumber {
			blockNumber, err = stages.GetStageProgress(dbReader, stages.Execution)
			if err != nil {
				return 0, common.Hash{}, fmt.Errorf("getting latest block number: %v", err)
			}

		} else if number == rpc.EarliestBlockNumber {
			blockNumber = 0

		} else if number == rpc.PendingBlockNumber {
			return 0, common.Hash{}, fmt.Errorf("pending blocks are not supported")

		} else {
			blockNumber = uint64(number.Int64())
		}
		hash, err = rawdb.ReadCanonicalHash(dbReader, blockNumber)
		if err != nil {
			return 0, common.Hash{}, err
		}
	} else {
		number := rawdb.ReadHeaderNumber(dbReader, hash)
		if number == nil {
			return 0, common.Hash{}, fmt.Errorf("block %x not found", hash)
		}
		blockNumber = *number

		ch, err := rawdb.ReadCanonicalHash(dbReader, blockNumber)
		if err != nil {
			return 0, common.Hash{}, err
		}
		if blockNrOrHash.RequireCanonical && ch != hash {
			return 0, common.Hash{}, fmt.Errorf("hash %q is not currently canonical", hash.String())
		}
	}
	return blockNumber, hash, nil
}

func GetAccount(tx ethdb.Database, blockNumber uint64, address common.Address) (*accounts.Account, error) {
	reader := adapter.NewStateReader(tx.(ethdb.HasTx).Tx(), blockNumber)
	return reader.ReadAccountData(address)
}
