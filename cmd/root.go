/*
Copyright © 2024 Joseph Zhao pandaski@outlook.com.au

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
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

func NewRootCmd(appFs afero.Fs) *cobra.Command {
	AppFs = appFs
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
	cmd.AddCommand(generateCmd)
	cmd.AddCommand(getCmd)
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

func init() {
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
