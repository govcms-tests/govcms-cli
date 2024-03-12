package cmd

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/govcms"
	"github.com/spf13/cobra"
)

// checkCmd represents the requirements command
var checkCmd = &cobra.Command{
	Use:    "check",
	Short:  "Check system requirements before running commands",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := govcms.CheckRequirements(); err != nil {
			// Customize the error message with the error on the next line
			return fmt.Errorf("System requirements check failed:\n%v", err)
		}
		fmt.Println("All system requirements met!")
		return nil
	},
}

func preRun(cmd *cobra.Command, args []string) error {
	// Check the requirements
	if err := govcms.CheckRequirements(); err != nil {
		// Customize the error message with the error on the next line
		return fmt.Errorf("System requirements check failed:\n%v", err)
	}
	return nil
}
