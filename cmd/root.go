package cmd

import (
	"github.com/polipopoliko/go-block/bootstrap"
)

func Execute() error {

	var (
		container = bootstrap.Init()
		node      = nodeCmd(container)
	)

	goblockCmd.AddCommand(
		wsCmd(container),
		node,
	)

	node.PersistentFlags().IntVarP(&port, "ports", "p", 9000, "running on ports")
	node.PersistentFlags().StringSliceVarP(&dial, "dial", "d", []string{}, "running on ports")
	node.PersistentFlags().IntVarP(&seed, "seed", "s", 0, "random seed")

	return goblockCmd.ExecuteContext(container.Context())
}
