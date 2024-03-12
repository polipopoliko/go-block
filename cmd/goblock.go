package cmd

import (
	"log"

	server "github.com/polipopoliko/go-block/api/infura"
	"github.com/polipopoliko/go-block/bootstrap"
	"github.com/spf13/cobra"
)

var (
	goblockCmd = func(container *bootstrap.Container) *cobra.Command {
		return &cobra.Command{
			PreRun: func(cmd *cobra.Command, args []string) {
				log.Println("running go block command")
			},
			Short: "running a golang based blockchain that connects to ethereum network",
			Use:   "goblock",

			Run: func(cmd *cobra.Command, args []string) {
				server.Serve(container)
			},
		}
	}
)
