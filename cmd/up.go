package cmd

import (
	"fmt"
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

		installPath = local.GetInstallPath(name)
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
