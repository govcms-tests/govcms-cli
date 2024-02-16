package cmd

import (
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise GovCMS CLI database",
	Long:  "Initialise GovCMS CLI database",
	Run: func(cmd *cobra.Command, args []string) {
		data.CreateInstallationTables()
	},
}

func init() {
	//RootCmd.AddCommand(initCmd)
}
