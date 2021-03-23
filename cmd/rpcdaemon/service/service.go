package service

import (
	"github.com/EVRICE/tgeth_alpha/cmd/rpcdaemon/cli"
	"github.com/EVRICE/tgeth_alpha/cmd/rpcdaemon/commands"
	"github.com/EVRICE/tgeth_alpha/core"
	"github.com/EVRICE/tgeth_alpha/ethdb"
	"github.com/EVRICE/tgeth_alpha/node"
)

func New(db ethdb.Database, ethereum core.Backend, stack *node.Node) {
	apis := commands.APIList(db, core.NewEthBackend(ethereum), nil, cli.Flags{API: []string{"eth", "debug"}}, nil)

	stack.RegisterAPIs(apis)
}
