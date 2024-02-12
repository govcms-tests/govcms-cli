package cmd

import (
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find [root_path]",
	Short: "Find all GovCMS installations from a given path",
	Long:  "Find all GovCMS installations from a given path. An 'installation' is defined to be any directory containing a 'govcms.info.yml' file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
