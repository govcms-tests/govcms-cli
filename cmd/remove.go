package cmd

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a local GovCMS installation",
	Long:  "Remove a local GovCMS installation",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range args {
			err := installationManager.DeleteInstallation(name)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
