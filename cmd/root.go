package cmd

import (
	"librenote/infrastructure/config"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath string

func rootCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "librenote",
		Short: "LibreNote is a note taking application",
		Long:  `LibreNote is a note taking application server`,
	}

	c.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "path to config file")

	c.AddCommand(newVersionCommand())
	c.AddCommand(migrateCommand())
	c.AddCommand(serveCommand())

	return c
}

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	err := config.Load(configPath)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// log config
	logrus.SetLevel(logrus.InfoLevel)

	if err := rootCommand().Execute(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}
