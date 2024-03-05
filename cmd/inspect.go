package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Clone, build, and launch a local copy of a GovCMS site",
	Long:  "Clone, build, and launch a local copy of a GovCMS site",
	Args:  cobra.ExactArgs(2),
	RunE:  inspect,
}

func inspect(cmd *cobra.Command, args []string) error {
	org := args[0]
	project := args[1]

	repoURL := "https://projects.govcms.gov.au/" + org + "/" + project
	err := installationManager.AddRaw(repoURL, appConfig.Workspace)
	if err != nil {
		fmt.Println("Error cloning repo")
		return err
	}

	LaunchSaas(project)

	return nil
}
