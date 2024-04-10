package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all locally installed GovCMS installations",
	Long:  "List all locally installed GovCMS installations",
	Run: func(cmd *cobra.Command, args []string) {
		installationManager.Sync()
		paths, _ := installationManager.GetAllPaths()

		// Render the table
		renderTable(paths)
	},
}

// Function to render a table with the provided paths
func renderTable(paths []string) {
	// Initialise a new table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredDark)

	// Define column headers
	t.SetTitle("Found the following local instances:")
	t.AppendHeader(table.Row{"#", "Path"})

	// Append rows to the table
	for i, path := range paths {
		// Add cells to the row with style for the first column
		row := table.Row{
			i + 1, // Set style for the first column
			path,
		}
		t.AppendRow(row)
	}

	// Render the table
	t.Render()
}
