package govcms

import (
	"fmt"

	"github.com/govcms-tests/govcms-cli/pkg/config"
	docker "github.com/govcms-tests/govcms-cli/pkg/utils"
)

// Pull the Docker images specified in the config
func Update() {
	// Access the DockerImages slice from the config package
	dockerImages := config.DockerImages

	// Iterate over the image list
	for _, image := range dockerImages {
		fmt.Printf("Pulling Docker image: %s\n", image)
		err := docker.DockerImagePull(image)
		if err != nil {
			fmt.Printf("Error pulling image %s: %v\n", image, err)
			continue // Continue to the next image on error
		}
		fmt.Printf("Successfully pulled image: %s\n", image)
	}
}
