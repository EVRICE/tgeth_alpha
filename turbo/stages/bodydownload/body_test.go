package bodydownload

import (
	"testing"

	"github.com/EVRICE/tgeth_alpha/ethdb"
)

func TestCreateBodyDownload(t *testing.T) {
	db := ethdb.NewMemDatabase()
	defer db.Close()
	bd := NewBodyDownload(100)
	if _, _, _, err := bd.UpdateFromDb(db); err != nil {
		t.Fatalf("update from db: %v", err)
	}
}
