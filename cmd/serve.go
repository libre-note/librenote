package cmd

import (
	"librenote/app/server"

	"github.com/spf13/cobra"
)

func serveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve serves the librenote api service",
		Long:  `Serve serves the librenote api service`,
		Run: func(cmd *cobra.Command, args []string) {
			server.Serve()
		},
	}
}
