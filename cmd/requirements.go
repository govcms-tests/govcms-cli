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

	"github.com/govcms-tests/govcms-cli/pkg/govcms"
	"github.com/spf13/cobra"
)

// requirementsCmd represents the requirements command
var requirementsCmd = &cobra.Command{
	Use:    "requirements",
	Short:  "Check system requirements before running commands",
	Hidden: true, // Hide this command from help messages
	Run: func(cmd *cobra.Command, args []string) {
		govcms.CheckRequirements()
	},
}

func init() {
	RootCmd.AddCommand(requirementsCmd)

	// Register the persistent pre-run function
	RootCmd.PersistentPreRunE = preRun
}

func preRun(cmd *cobra.Command, args []string) error {
	// Check the requirements
	if err := govcms.CheckRequirements(); err != nil {
		// Customize the error message with the error on the next line
		return fmt.Errorf("System requirements check failed:\n%v", err)
	}
	return nil
}
