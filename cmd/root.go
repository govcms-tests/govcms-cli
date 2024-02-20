package cmd

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/data"
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
var local data.LocalStorage

func NewRootCmd(appFs afero.Fs, localStorage data.LocalStorage) *cobra.Command {
	AppFs = appFs
	local = localStorage
	cmd := &cobra.Command{
		Use:     "govcms",
		Short:   "Lift the GovCMS local development",
		Long:    "Lift the GovCMS local development",
		Version: "v0.1.0 -- HEAD",
	}
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.govcms-cli.yaml)")
	cmd.SetVersionTemplate("GovCMS CLI version " + cmd.Version + "\n")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(cleanupCmd)
	cmd.AddCommand(distributionCmd)
	cmd.AddCommand(findCmd)
	//cmd.AddCommand(generateCmd)
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(guiCmd)
	cmd.AddCommand(initCmd)
	cmd.AddCommand(issueCmd)
	cmd.AddCommand(listCmd)
	cmd.AddCommand(requirementsCmd)
	cmd.AddCommand(testCmd)
	cmd.AddCommand(upCmd)
	cmd.AddCommand(downCmd)
	cmd.AddCommand(updateCmd)
	cmd.AddCommand(versionCmd)

	// Register the persistent pre-run function
	cmd.PersistentPreRunE = preRun

	return cmd
}

//func init() {
//}

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
