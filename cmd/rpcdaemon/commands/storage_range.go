package commands

import (
	"fmt"

	"github.com/holiman/uint256"
	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/turbo/adapter"
)

// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage StorageMap   `json:"storage"`
	NextKey *common.Hash `json:"nextKey"` // nil if Storage includes the last key in the trie.
}

// StorageMap a map from storage locations to StorageEntry items
type StorageMap map[common.Hash]StorageEntry

// StorageEntry an entry in storage of the account
type StorageEntry struct {
	Key   *common.Hash `json:"key"`
	Value common.Hash  `json:"value"`
}

func StorageRangeAt(stateReader *adapter.StateReader, contractAddress common.Address, start []byte, maxResult int) (StorageRangeResult, error) {
	result := StorageRangeResult{Storage: StorageMap{}}
	resultCount := 0

	if err := stateReader.ForEachStorage(contractAddress, common.BytesToHash(start), func(key, seckey common.Hash, value uint256.Int) bool {
		if resultCount < maxResult {
			result.Storage[seckey] = StorageEntry{Key: &key, Value: value.Bytes32()}
		} else {
			result.NextKey = &key
		}
		resultCount++
		return resultCount <= maxResult
	}, maxResult+1); err != nil {
		return StorageRangeResult{}, fmt.Errorf("error walking over storage: %v", err)
	}
	return result, nil
}
