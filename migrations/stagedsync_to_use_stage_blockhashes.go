package migrations

import (
	"github.com/EVRICE/tgeth_alpha/common/etl"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync/stages"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

var stagedsyncToUseStageBlockhashes = Migration{
	Name: "stagedsync_to_use_stage_blockhashes",
	Up: func(db ethdb.Database, tmpdir string, progress []byte, OnLoadCommit etl.LoadCommitHandler) error {

		var stageProgress uint64
		var err error
		if stageProgress, err = stages.GetStageProgress(db, stages.Headers); err != nil {
			return err
		}

		if err = stages.SaveStageProgress(db, stages.BlockHashes, stageProgress); err != nil {
			return err
		}

		if err = OnLoadCommit(db, nil, true); err != nil {
			return err
		}

		return nil
	},
}

var unwindStagedsyncToUseStageBlockhashes = Migration{
	Name: "unwind_stagedsync_to_use_stage_blockhashes",
	Up: func(db ethdb.Database, tmpdir string, progress []byte, OnLoadCommit etl.LoadCommitHandler) error {

		var stageProgress uint64
		var err error
		if stageProgress, err = stages.GetStageUnwind(db, stages.Headers); err != nil {
			return err
		}

		if err = stages.SaveStageUnwind(db, stages.BlockHashes, stageProgress); err != nil {
			return err
		}

		if err = OnLoadCommit(db, nil, true); err != nil {
			return err
		}

		return nil
	},
}
