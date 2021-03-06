package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/EVRICE/tgeth_alpha/cmd/utils"
	"github.com/EVRICE/tgeth_alpha/internal/debug"
	"github.com/EVRICE/tgeth_alpha/log"
	"github.com/spf13/cobra"
)

func init() {
	utils.CobraFlags(rootCmd, append(debug.Flags, utils.MetricFlags...))

}
func Execute() {
	if err := rootCmd.ExecuteContext(rootContext()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-ch:
			log.Info("Got interrupt, shutting down...")
		case <-ctx.Done():
		}

		cancel()
	}()
	return ctx
}

var (
	chaindata    string
	snapshotFile string
	block        uint64
	snapshotDir  string
	snapshotMode string
)

var rootCmd = &cobra.Command{
	Use:   "generate_snapshot",
	Short: "generate snapshot",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := debug.SetupCobra(cmd); err != nil {
			panic(err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		debug.Exit()
	},
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func withBlock(cmd *cobra.Command) {
	cmd.Flags().Uint64Var(&block, "block", 1, "specifies a block number for operation")
}
func withSnapshotData(cmd *cobra.Command) {
	cmd.Flags().StringVar(&snapshotMode, "snapshotMode", "", "set of snapshots to use")
	cmd.Flags().StringVar(&snapshotDir, "snapshotDir", "", "snapshot dir")
}

func withChaindata(cmd *cobra.Command) {
	cmd.Flags().StringVar(&chaindata, "chaindata", "chaindata", "path to the chaindata file used as input to analysis")
	must(cmd.MarkFlagFilename("chaindata", ""))
}

func withSnapshotFile(cmd *cobra.Command) {
	cmd.Flags().StringVar(&snapshotFile, "snapshot", "", "path where to write the snapshot file")
}
