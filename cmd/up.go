package cmd

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/data"
	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

type saveOutput struct {
	savedOutput []byte
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

var upCmd = &cobra.Command{
	Use:   "up [resource]",
	Short: "Launch docker container",
	Run: func(cmd *cobra.Command, args []string) {
		var so saveOutput
		var installPath string
		var name string

		if len(args) < 1 {
			installPath = "govcms/distribution"
		} else {
			name = args[0]
		}

		installPath = data.GetInstallPath(name)
		fmt.Println("Attempting to launch site located at " + installPath)

		command := exec.Command("/bin/sh", "-c", "docker compose up -d")
		command.Dir = installPath

		command.Stdin = os.Stdin
		command.Stdout = &so
		command.Stderr = os.Stderr

		_ = command.Run()
		fmt.Printf("%s", so.savedOutput)

		fmt.Println("\nLocal server has started at", termlink.Link("http://localhost:8888", "http://localhost:8888"))
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop docker container",
	Run: func(cmd *cobra.Command, arg []string) {
		var so saveOutput

		command := exec.Command("/bin/sh", "-c", "docker compose down")
		command.Dir = "govcms/distribution"

		command.Stdin = os.Stdin
		command.Stdout = &so
		command.Stderr = os.Stderr

		_ = command.Run()
		fmt.Printf("%s", so.savedOutput)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}
