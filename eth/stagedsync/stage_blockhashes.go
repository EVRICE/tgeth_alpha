package stagedsync

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/common/dbutils"
	"github.com/EVRICE/tgeth_alpha/common/etl"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync/stages"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

func extractHeaders(k []byte, v []byte, next etl.ExtractNextFunc) error {
	// We only want to extract entries composed by Block Number + Header Hash
	if len(k) != 40 {
		return nil
	}
	return next(k, common.CopyBytes(k[8:]), common.CopyBytes(k[:8]))
}

func SpawnBlockHashStage(s *StageState, db ethdb.Database, tmpdir string, quit <-chan struct{}) error {
	var tx ethdb.DbWithPendingMutations
	var useExternalTx bool
	if hasTx, ok := db.(ethdb.HasTx); ok && hasTx.Tx() != nil {
		tx = db.(ethdb.DbWithPendingMutations)
		useExternalTx = true
	} else {
		var err error
		tx, err = db.Begin(context.Background(), ethdb.RW)
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}
	headNumber, err := stages.GetStageProgress(tx, stages.Headers)
	if err != nil {
		return fmt.Errorf("getting headers progress: %w", err)
	}
	headHash := rawdb.ReadHeaderByNumber(tx, headNumber).Hash()
	if s.BlockNumber == headNumber {
		s.Done()
		return nil
	}

	startKey := make([]byte, 8)
	binary.BigEndian.PutUint64(startKey, s.BlockNumber)
	endKey := dbutils.HeaderKey(headNumber, headHash) // Make sure we stop at head

	logPrefix := s.state.LogPrefix()
	if err := etl.Transform(
		logPrefix,
		tx,
		dbutils.HeaderPrefix,
		dbutils.HeaderNumberPrefix,
		tmpdir,
		extractHeaders,
		etl.IdentityLoadFunc,
		etl.TransformArgs{
			ExtractStartKey: startKey,
			ExtractEndKey:   endKey,
			Quit:            quit,
		},
	); err != nil {
		return err
	}
	if err := s.DoneAndUpdate(tx, headNumber); err != nil {
		return err
	}
	if !useExternalTx {
		if _, err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
