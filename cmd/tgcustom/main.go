package main

import (
	"fmt"
	"os"

	"github.com/EVRICE/tgeth_alpha/common/dbutils"
	"github.com/EVRICE/tgeth_alpha/core/state"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync"
	"github.com/EVRICE/tgeth_alpha/eth/stagedsync/stages"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/log"
	"github.com/EVRICE/tgeth_alpha/turbo/node"

	turbocli "github.com/EVRICE/tgeth_alpha/turbo/cli"

	"github.com/urfave/cli"
)

// defining a custom command-line flag, a string
var flag = cli.StringFlag{
	Name:  "custom-stage-greeting",
	Value: "default-value",
}

// defining a custom bucket name
const (
	customBucketName = "ch.torquem.demo.tgcustom.CUSTOM_BUCKET"
)

// the regular main function
func main() {
	// initializing turbo-geth application here and providing our custom flag
	app := turbocli.MakeApp(runTurboGeth,
		append(turbocli.DefaultFlags, flag), // always use DefaultFlags, but add a new one in the end.
	)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func syncStages(ctx *cli.Context) stagedsync.StageBuilders {
	return append(
		stagedsync.DefaultStages(), // adding all default stages
		stagedsync.StageBuilder{ // adding our custom stage
			ID: stages.SyncStage("ch.torquem.demo.tgcustom.CUSTOM_STAGE"),
			Build: func(world stagedsync.StageParameters) *stagedsync.Stage {
				return &stagedsync.Stage{
					ID:          stages.SyncStage("ch.torquem.demo.tgcustom.CUSTOM_STAGE"),
					Description: "Custom Stage",
					ExecFunc: func(s *stagedsync.StageState, _ stagedsync.Unwinder) error {
						fmt.Println("hello from the custom stage", ctx.String(flag.Name))
						val, err := world.TX.Get(customBucketName, []byte("test"))
						fmt.Println("val", string(val), "err", err)
						if err := world.TX.Put(customBucketName, []byte("test"), []byte(ctx.String(flag.Name))); err != nil {
							return err
						}
						s.Done()
						return nil
					},
					UnwindFunc: func(u *stagedsync.UnwindState, s *stagedsync.StageState) error {
						fmt.Println("hello from the custom stage unwind", ctx.String(flag.Name))
						if err := world.TX.Delete(customBucketName, []byte("test"), nil); err != nil {
							return err
						}
						return u.Done(world.TX)
					},
				}
			},
		},
	)
}

// turbo-geth main function
func runTurboGeth(ctx *cli.Context) {
	// creating a staged sync with our new stage
	sync := stagedsync.New(
		syncStages(ctx),
		stagedsync.DefaultUnwindOrder(),
		stagedsync.OptionalParameters{
			StateReaderBuilder: func(db ethdb.Database) state.StateReader {
				// put your custom caching code here
				return state.NewPlainStateReader(db)
			},
			StateWriterBuilder: func(db ethdb.Database, changeSetsDB ethdb.Database, blockNumber uint64) state.WriterWithChangeSets {
				// put your custom cache update code here
				return state.NewPlainStateWriter(db, changeSetsDB, blockNumber)
			},
		},
	)

	// running a node and initializing a custom bucket with all default settings
	tg := node.New(ctx, sync, node.Params{
		CustomBuckets: map[string]dbutils.BucketConfigItem{
			customBucketName: {},
		},
	})

	err := tg.Serve()

	if err != nil {
		log.Error("error while serving a turbo-geth node", "err", err)
	}
}
