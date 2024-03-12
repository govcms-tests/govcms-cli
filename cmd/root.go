package cmd

import (
	database "database-sqlc"
	"database/sql"
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

// var RootCmd *cobra.Command
var cfgFile string
var appConfig settings.Config // Rename config to appConfig
var AppFs afero.Fs
var installationManager *database.InstallationManager

func NewRootCmd(appFs afero.Fs, db *sql.DB) *cobra.Command {
	AppFs = appFs
	installationManager = database.NewInstallationManager(db, appFs)
	cmd := &cobra.Command{
		Use:     "govcms",
		Short:   "A CLI tool to help with the development and maintenance of the GovCMS platform",
		Long:    "A CLI tool to help with the development and maintenance of the GovCMS platform",
		Version: "0.2.1",
	}
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Set config file, default is $HOME/.govcms-cli.yaml")
	cmd.SetVersionTemplate("GovCMS CLI version " + cmd.Version + "\n")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(removeCmd)
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(upCmd)
	cmd.AddCommand(downCmd)
	cmd.AddCommand(checkCmd)
	cmd.AddCommand(listCmd)

	// Register the persistent pre-run function
	cmd.PersistentPreRunE = preRun

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Load configuration settings from settings.go
	var err error
	appConfig, err = settings.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".govcms-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".govcms-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("Using config file:", viper.ConfigFileUsed())
	}
}
