package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewGetCmd() *cobra.Command {
	Reader = os.Stdin

	cmd := &cobra.Command{
		Use:   "get [resource] [name]",
		Short: "get a GovCMS distribution, saas, or paas site",
		Long:  "get a GovCMS distribution, saas, or paas site.",
		RunE:  runGetCommand,
	}
	cmd.Flags().IntP("pr", "p", 0, "Github PR number")
	cmd.Flags().StringP("branch", "b", "", "Git branch name")

	return cmd
}

func runGetCommand(cmd *cobra.Command, args []string) error {
	if hasBothFlags(cmd) {
		return fmt.Errorf("cannot specify both --pr and --branch flags together")
	}
	if len(args) < 2 {
		err := getGovcmsWithPrompt()
		return err
	}
	err := getGovcmsWithoutPrompt(cmd, args)
	return err
}

func hasBothFlags(cmd *cobra.Command) bool {
	prNumber, _ := cmd.Flags().GetInt("pr")
	branch, _ := cmd.Flags().GetString("branch")
	return prNumber != 0 && branch != ""
}

func getGovcmsWithPrompt() error {
	name, govcmsType := getDetailsFromPrompt()
	err := CloneGovCMS(name, govcmsType, "", 0)
	return err
}

func getDetailsFromPrompt() (string, string) {
	name := getNameFromPrompt()
	govcmsType := getTypeFromPrompt()
	return name, govcmsType
}

func getNameFromPrompt() string {
	prompt := promptui.Prompt{
		Label: "What would you like to call this installation?",
		Stdin: Reader,
	}
	name, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	fmt.Printf("You chose %q\n", name)
	return name
}

func getTypeFromPrompt() string {
	prompt := promptui.Select{
		Label: "Which type would you like to install?",
		Items: []string{"Distribution", "SaaS", "PaaS", "Lagoon", "Tests", "Scaffold-Tooling"},
		Stdin: Reader,
	}
	_, govcmsType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	fmt.Printf("You choose %q\n", govcmsType)
	return strings.ToLower(govcmsType)
}

func getGovcmsWithoutPrompt(cmd *cobra.Command, args []string) error {
	govcmsType := args[0]
	name := args[1]
	prNumber, _ := cmd.Flags().GetInt("pr")
	branchName, _ := cmd.Flags().GetString("branch")
	// Call the generate function from the govcms package
	err := CloneGovCMS(name, govcmsType, branchName, prNumber)
	return err
}
