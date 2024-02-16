/*
Copyright © 2024 Joseph Zhao pandaski@outlook.com.au

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/govcms-tests/govcms-cli/pkg/govcms"
	"github.com/govcms-tests/govcms-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// distributionCmd represents the distribution command
var distributionCmd = &cobra.Command{
	Use:   "distribution [task]",
	Short: "Manage GovCMS distribution for local development",
	Long:  "Manage the GovCMS distribution for local development.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		// Check if the specified path is a valid GovCMS distribution
		if _, err := os.Stat(path + "/composer.json"); os.IsNotExist(err) {
			fmt.Println("Error: Invalid GovCMS distribution path.")
			return
		}

		// If no arguments provided, display basic information about the current folder
		if len(args) == 0 {
			fmt.Println("Path: ", path)

			// Call the function to calculate file info
			fileCount, totalSize, err := utils.CalculateFileInfo(path)
			if err != nil {
				fmt.Println("Error calculating file info:", err)
				return
			}
			fmt.Printf("Total number of files: %d\n", fileCount)
			fmt.Printf("Total size of all files: %.2f MB\n", float64(totalSize)/1048576) // Convert bytes to megabytes
			return
		}

		task := args[0]

		// Only execute the setup steps if arguments are provided
		if task == "up" || task == "review" || task == "test" {
			govcms.CopyDistributionFiles(path)
			govcms.StartContainer(path)
		}

		switch task {
		case "up":
			// Already executed setup steps
		case "review":
			fmt.Println("Feature not implemented yet.")
		case "test":
			fmt.Println("Feature not implemented yet.")
		case "cleanup":
			govcms.StopContainer(path)
			fmt.Println("Docker container stopped and cleaned up.")
		default:
			fmt.Println("Error: invalid subcommand.")
		}
	},
}

func init() {
	//RootCmd.AddCommand(distributionCmd)
	distributionCmd.Flags().StringP("path", "p", ".", "Path to GovCMS distribution folder")
}
