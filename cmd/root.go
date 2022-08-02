package cmd

import (
	"fmt"
	"librenote/app"
	"librenote/infrastructure/config"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// nolint:gochecknoglobals
var (
	configPath string
	rootCmd    = &cobra.Command{
		Use:     "librenote",
		Short:   "LibreNote is a note taking application",
		Long:    `LibreNote is a note taking application server`,
		Version: app.Version,
	}
)

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate(fmt.Sprintf("LibreNote Core\nVersion: %s\nBuild time: %s\n", app.Version, app.BuildTime))
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "path to config file")
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

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}
