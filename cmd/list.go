package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all locally installed GovCMS installations",
	Long:  "List all locally installed GovCMS installations",
	Run: func(cmd *cobra.Command, args []string) {
		installationManager.Sync()
		fmt.Println("Found the following local instances:")
		fmt.Println("============================")
		paths, _ := installationManager.GetAllPaths()
		fmt.Println(strings.Join(paths, "\n"))
	},
}
