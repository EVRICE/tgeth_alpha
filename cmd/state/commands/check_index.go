package commands

import (
	"github.com/EVRICE/tgeth_alpha/cmd/state/verify"
	"github.com/EVRICE/tgeth_alpha/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	withChaindata(checkIndexCMD)
	withIndexBucket(checkIndexCMD)
	withCSBucket(checkIndexCMD)
	rootCmd.AddCommand(checkIndexCMD)
}

var checkIndexCMD = &cobra.Command{
	Use:   "checkIndex",
	Short: "Index checker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify.CheckIndex(utils.RootContext(), chaindata, changeSetBucket, indexBucket)
	},
}
