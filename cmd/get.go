package cmd

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/govcms"
	"github.com/manifoldco/promptui"
	"strings"

	"github.com/spf13/cobra"
)

var cmd cobra.Command

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [resource] [name]",
	Short: "get a GovCMS distribution, saas, or paas site",
	Long:  "get a GovCMS distribution, saas, or paas site.",
	//Args:      cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	//ValidArgs: []string{"distribution", "saas", "paas"},
	Run: func(cmd *cobra.Command, args []string) {
		if hasBothFlags() {
			fmt.Println("Error: Cannot specify both --pr and --branch flags together.")
			return
		}
		if len(args) < 2 {
			getGovcmsWithPrompt()
			return
		}
		getGovcmsWithoutPrompt(args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func hasBothFlags() bool {
	return cmd.Flags().Changed("pr") && cmd.Flags().Changed("branch")
}

func getGovcmsWithPrompt() {
	name, govcmsType := getDetailsFromPrompt()
	err := govcms.Generate(name, govcmsType, 0, "")
	if err != nil {
		fmt.Printf("Error generating %s: %v\n", govcmsType, err)
		return
	}
}

func getDetailsFromPrompt() (string, string) {
	name := getNameFromPrompt()
	govcmsType := getTypeFromPrompt()
	return name, govcmsType
}

func getNameFromPrompt() string {
	prompt := promptui.Prompt{
		Label: "What would you like to call this installation?",
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
	}
	_, govcmsType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	fmt.Printf("You choose %q\n", govcmsType)
	return strings.ToLower(govcmsType)
}

func getGovcmsWithoutPrompt(args []string) {
	govcmsType := args[0]
	name := args[1]
	prNumber, _ := cmd.Flags().GetInt("pr")
	branchName, _ := cmd.Flags().GetString("branch")
	// Call the generate function from the govcms package
	err := govcms.Generate(name, govcmsType, prNumber, branchName)
	if err != nil {
		fmt.Printf("Error generating %s: %v\n", govcmsType, err)
		return
	}
}
