package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCommand creates a new root cli command instance
func RootCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "librenote",
		Short: "LibreNote is a note taking application",
		Long:  `LibreNote is a note taking application server`,
	}

	c.AddCommand(NewVersionCommand())

	return c
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// log config
	logrus.SetLevel(logrus.InfoLevel)

	if err := RootCommand().Execute(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}
