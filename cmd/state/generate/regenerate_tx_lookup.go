package generate

import (
	"os"
	"os/signal"
	"time"

	"github.com/EVRICE/tgeth_alpha/common/dbutils"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync/stages"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/log"
)

func RegenerateTxLookup(chaindata string) error {
	db := ethdb.MustOpen(chaindata)
	defer db.Close()
	if err := db.ClearBuckets(dbutils.TxLookupPrefix); err != nil {
		return err
	}
	startTime := time.Now()
	ch := make(chan os.Signal, 1)
	quitCh := make(chan struct{})
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		close(quitCh)
	}()

	lastExecutedBlock, err := stages.GetStageProgress(db, stages.Execution)
	if err != nil {
		//There could be headers without block in the end
		log.Error("Cant get last executed block", "err", err)
	}
	log.Info("TxLookup generation started", "start time", startTime)
	err = stagedsync.TxLookupTransform("txlookup", db, dbutils.HeaderHashKey(0), dbutils.HeaderHashKey(lastExecutedBlock), quitCh, os.TempDir())
	if err != nil {
		return err
	}
	log.Info("TxLookup index is successfully regenerated", "it took", time.Since(startTime))
	return nil
}
