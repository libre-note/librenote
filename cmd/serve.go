package cmd

import (
	"librenote/app/server"

	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(serveCmd)
}

// nolint:gochecknoglobals
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve serves the librenote api service",
	Long:  `Serve serves the librenote api service`,
	Run: func(cmd *cobra.Command, args []string) {
		serverReady := make(chan bool)
		s := server.Server{ServerReady: serverReady}
		s.Serve()
	},
}
