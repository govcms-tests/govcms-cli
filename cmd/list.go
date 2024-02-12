package cmd

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/data"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all locally installed GovCMS instances",
	Long:  "List all locally installed GovCMS instances",
	Run: func(cmd *cobra.Command, args []string) {
		data.SyncInstallations()
		fmt.Println(strings.Join(data.GetListOfPaths(), "\n"))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
