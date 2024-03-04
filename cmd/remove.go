package cmd

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Navigate to an installation's local file directory",
	Long:  "Navigate to an installation's local file directory",
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
