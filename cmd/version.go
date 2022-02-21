package cmd

import (
	"fmt"
	"librenote/app"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Print the version number of LibreNote librenote",
		Run:   printVersion,
	}
}

func printVersion(_ *cobra.Command, _ []string) {
	fmt.Println("LibreNote Core")
	fmt.Printf("Version: %s\n", app.Version)
	fmt.Printf("Build time: %s\n", app.BuildTime)
}
