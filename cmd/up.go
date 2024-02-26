package cmd

import (
	_ "embed"
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/utils"
	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Up(cmd, args)
	},
}

func Up(cmd *cobra.Command, args []string) {
	var so saveOutput
	var installPath string
	name := args[0]

	// Prepare command execution
	installPath = local.GetInstallPath(name)
	_ = setRandomPort(installPath)
	command := exec.Command("/bin/sh", "-c", "docker compose up -d")
	command.Dir = installPath
	command.Stdin = os.Stdin
	command.Stdout = &so
	command.Stderr = os.Stderr

	// Execute up command
	_ = command.Run()
	fmt.Printf("%s", so.savedOutput)

	ip := utils.GetContainerIpByName(name)
	fmt.Println("\nLocal server has started at", termlink.Link("http://"+ip, "http://"+ip))
}

func setRandomPort(path string) error {
	port := getFreePort(49152, 65535)
	replaceCmd := "sed -i \"\" -E \"s/[0-9]+:80/" + strconv.Itoa(port) + ":80/g\" docker-compose.yml"

	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(replaceCmd)
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run script: %v", err)
	}
	return nil
}

func getFreePort(min int, max int) int {
	// Select random ports until one is found that is free
	for {
		randPort := rand.Intn(max-min+1) + min
		if free, _ := isPortFree(randPort); free {
			return randPort
		}
	}
}

// Check if a port is available
func isPortFree(port int) (status bool, err error) {
	// Concatenate a colon and the port
	host := ":" + strconv.Itoa(port)
	// Try to create a server with the port
	server, err := net.Listen("tcp", host)
	// if it fails then the port is likely taken
	if err != nil {
		return false, err
	}
	// close the server
	server.Close()
	// we successfully used and closed the port
	// so it's now available to be used again
	return true, nil
}
