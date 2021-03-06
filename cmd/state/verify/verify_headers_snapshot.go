package verify

import (
	"context"
	"errors"

	"github.com/ledgerwatch/lmdb-go/lmdb"
	"github.com/EVRICE/tgeth_alpha/common/dbutils"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/log"
	"github.com/EVRICE/tgeth_alpha/rlp"
)

func HeadersSnapshot(snapshotPath string) error {
	snKV := ethdb.NewLMDB().Path(snapshotPath).Flags(func(flags uint) uint { return flags | lmdb.Readonly }).WithBucketsConfig(func(defaultBuckets dbutils.BucketsCfg) dbutils.BucketsCfg {
		return dbutils.BucketsCfg{
			dbutils.HeaderPrefix:              dbutils.BucketConfigItem{},
			dbutils.HeadersSnapshotInfoBucket: dbutils.BucketConfigItem{},
		}
	}).MustOpen()
	var prevHeader *types.Header
	err := snKV.View(context.Background(), func(tx ethdb.Tx) error {
		c := tx.Cursor(dbutils.HeaderPrefix)
		k, v, innerErr := c.First()
		for {
			if len(k) == 0 && len(v) == 0 {
				break
			}
			if innerErr != nil {
				return innerErr
			}

			header := new(types.Header)
			innerErr := rlp.DecodeBytes(v, header)
			if innerErr != nil {
				return innerErr
			}

			if prevHeader != nil {
				if prevHeader.Number.Uint64()+1 != header.Number.Uint64() {
					log.Error("invalid header number", "p", prevHeader.Number.Uint64(), "c", header.Number.Uint64())
					return errors.New("invalid header number")
				}
				if prevHeader.Hash() != header.ParentHash {
					log.Error("invalid parent hash", "p", prevHeader.Hash(), "c", header.ParentHash)
					return errors.New("invalid parent hash")
				}
			}
			k, v, innerErr = c.Next() //nolint
			prevHeader = header
		}
		return nil
	})
	return err
}
