package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop docker container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var so saveOutput
		name := args[0]

		installPath, err := installationManager.GetPath(name)
		if err != nil {
			return
		}

		command := exec.Command("/bin/sh", "-c", "docker compose down")
		command.Dir = installPath

		command.Stdin = os.Stdin
		command.Stdout = &so
		command.Stderr = os.Stderr

		_ = command.Run()
		fmt.Printf("%s", so.savedOutput)
	},
}
