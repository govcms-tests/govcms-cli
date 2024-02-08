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

	"github.com/govcms-tests/govcms-cli/utils/docker"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pulls Docker images and recreates the containers",
	Long:  "Pulls Docker images and recreates the containers.",
	Run: func(cmd *cobra.Command, args []string) {
		// List of Docker image names to pull
		imageList := []string{"govcmsextras/dnsmasq", "adminer", "mariadb:lts", "mailhog/mailhog"}
		// Iterate over the image list
		for _, imageName := range imageList {
			err := docker.DockerImagePull(imageName)
			if err != nil {
				fmt.Printf("Error pulling image %s: %v\n", imageName, err)
				continue // Continue to the next image on error
			}
			fmt.Printf("Successfully pulled image: %s\n", imageName)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
