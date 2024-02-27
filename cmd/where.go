package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Navigate to an installation's local file directory",
	Long:  "Navigate to an installation's local file directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installName := args[0]
		installPath := local.GetInstallPath(installName)
		cmd.SetOut(os.Stdout)
		cmd.Println(installPath)

	},
}
