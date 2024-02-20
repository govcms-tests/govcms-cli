package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all locally installed GovCMS instances",
	Long:  "List all locally installed GovCMS instances",
	Run: func(cmd *cobra.Command, args []string) {
		local.SyncInstallations()
		fmt.Println("Found the following local instances:")
		fmt.Println("============================")
		fmt.Println(strings.Join(local.GetListOfPaths(), "\n"))
	},
}
