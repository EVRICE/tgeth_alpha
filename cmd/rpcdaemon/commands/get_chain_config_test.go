package commands

import (
	"testing"

	"github.com/EVRICE/tgeth_alpha/core"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

func TestGetChainConfig(t *testing.T) {
	db := ethdb.NewMemDatabase()
	defer db.Close()
	config, _, _, err := core.SetupGenesisBlock(db, core.DefaultGenesisBlock(), false /* history */, false /* overwrite */)
	if err != nil {
		t.Fatalf("setting up genensis block: %v", err)
	}

	api := (&BaseAPI{})
	config1, err1 := api.chainConfig(db)
	if err1 != nil {
		t.Fatalf("reading chain config: %v", err1)
	}
	if config.String() != config1.String() {
		t.Fatalf("read different config: %s, expected %s", config1.String(), config.String())
	}
	config2, err2 := api.chainConfig(db)
	if err2 != nil {
		t.Fatalf("reading chain config: %v", err2)
	}
	if config.String() != config2.String() {
		t.Fatalf("read different config: %s, expected %s", config2.String(), config.String())
	}
}
