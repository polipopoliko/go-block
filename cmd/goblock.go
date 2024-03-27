package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	goblockCmd = &cobra.Command{
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Println("running go block command")
		},
		Short: "running a golang based blockchain that connects to ethereum network",
		Use:   "goblock",
	}
)
