package cmd

import (
	"fmt"
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
	Use:   "up",
	Short: "Launch docker container",
	Run: func(cmd *cobra.Command, arg []string) {
		var so saveOutput

		command := exec.Command("/bin/sh", "-c", "docker compose up -d")
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
}
