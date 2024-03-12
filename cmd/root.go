package cmd

import (
	"github.com/polipopoliko/go-block/bootstrap"
	"github.com/spf13/cobra"
)

func Execute() error {

	var (
		container = bootstrap.Init()
		root      = &cobra.Command{}
		node      = nodeCmd(container)
	)

	root.AddCommand(
		goblockCmd(container),
		node,
	)

	node.PersistentFlags().IntVarP(&port, "ports", "p", 9000, "running on ports")
	node.PersistentFlags().StringSliceVarP(&dial, "dial", "d", []string{}, "running on ports")
	node.PersistentFlags().IntVarP(&seed, "seed", "s", 0, "random seed")

	return root.ExecuteContext(container.Context())
}
