package headerdownload

import (
	"context"
	"math/big"
	"testing"

	"github.com/EVRICE/tgeth_alpha/common"
	"github.com/EVRICE/tgeth_alpha/core/rawdb"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

func TestInserter1(t *testing.T) {
	db := ethdb.NewMemDatabase()
	defer db.Close()
	// Set up parent difficulty
	if err := rawdb.WriteTd(db, common.Hash{}, 4, big.NewInt(0)); err != nil {
		t.Fatalf("write parent diff: %v", err)
	}
	tx, err := db.Begin(context.Background(), ethdb.RW)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	defer tx.Rollback()
	batch := tx.NewBatch()
	hi := NewHeaderInserter("headers", tx, batch, big.NewInt(0), 0)
	if err := hi.FeedHeader(&types.Header{Number: big.NewInt(5), Difficulty: big.NewInt(1)}, 5); err != nil {
		t.Errorf("feed empty header: %v", err)
	}
}
