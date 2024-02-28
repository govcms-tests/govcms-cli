package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Navigate to an installation's local file directory",
	Long:  "Navigate to an installation's local file directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range args {
			err := removeInstallationByName(name)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func removeInstallationByName(name string) error {
	installPath, err := local.GetInstallPath(name)
	if err != nil {
		fmt.Printf("Unable to find installation '%s'\n", name)
		return err
	}

	err = local.RemoveInstallFromName(name)
	if err != nil {
		return err
	}

	err = AppFs.RemoveAll(installPath)
	if err != nil {
		return err
	}
	fmt.Printf("Removed '%s' successfully\n", name)
	return nil
}
