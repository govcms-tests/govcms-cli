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
