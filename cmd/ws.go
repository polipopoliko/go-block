package cmd

import (
	"log"

	"github.com/polipopoliko/go-block/bootstrap"
	"github.com/spf13/cobra"
)

var (
	wsCmd = func(container *bootstrap.Container) *cobra.Command {
		return &cobra.Command{
			PreRun: func(cmd *cobra.Command, args []string) {
				log.Println("running ws command")
			},
			Short: "running a golang based blockchain that connects to ethereum network",
			Use:   "ws",

			Run: func(cmd *cobra.Command, args []string) {
			},
		}
	}
)
