package main

import (
	"fmt"
	"github.com/EVRICE/tgeth_alpha/cmd/integration/commands"
	"github.com/EVRICE/tgeth_alpha/cmd/utils"
	"os"
)

func main() {
	rootCmd := commands.RootCommand()

	if err := rootCmd.ExecuteContext(utils.RootContext()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
