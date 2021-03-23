package main

import (
	"os"

	"github.com/EVRICE/tgeth_alpha/cmd/rpcdaemon/cli"
	"github.com/EVRICE/tgeth_alpha/cmd/rpcdaemon/commands"
	"github.com/EVRICE/tgeth_alpha/cmd/rpcdaemon/filters"
	"github.com/EVRICE/tgeth_alpha/cmd/utils"
	"github.com/EVRICE/tgeth_alpha/common/fdlimit"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/log"
	"github.com/spf13/cobra"
)

func main() {
	raiseFdLimit()
	cmd, cfg := cli.RootCommand()
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		db, backend, err := cli.OpenDB(*cfg)
		if err != nil {
			log.Error("Could not connect to DB", "error", err)
			return nil
		}
		defer db.Close()

		var ff *filters.Filters
		if backend != nil {
			ff = filters.New(backend)
		} else {
			log.Info("filters are not supported in chaindata mode")
		}

		return cli.StartRpcServer(cmd.Context(), *cfg, commands.APIList(ethdb.NewObjectDatabase(db), backend, ff, *cfg, nil))
	}

	if err := cmd.ExecuteContext(utils.RootContext()); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

// raiseFdLimit raises out the number of allowed file handles per process
func raiseFdLimit() {
	limit, err := fdlimit.Maximum()
	if err != nil {
		log.Error("Failed to retrieve file descriptor allowance", "error", err)
		return
	}
	if _, err = fdlimit.Raise(uint64(limit)); err != nil {
		log.Error("Failed to raise file descriptor allowance", "error", err)
	}
}
