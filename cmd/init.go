package cmd

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise GovCMS CLI database",
	Long:  "Initialise GovCMS CLI database",
	Run: func(cmd *cobra.Command, args []string) {
		local.CreateInstallationTables()
	},
}

func init() {
	//RootCmd.AddCommand(initCmd)
}
