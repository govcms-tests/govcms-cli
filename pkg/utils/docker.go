package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// list docker image(s)
func DockerImageList() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		panic(err)
	}

	// list all Docker images
	imageList, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	// print the list of Docker images
	for _, image := range imageList {
		fmt.Println(image.ID, image.RepoTags)
	}
}

// Pulls a Docker image for a GovCMS service and prints the output in a nicer format.
func DockerImagePull(imageName string) error {
	// Create a Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Pull the Docker image
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	// Parse and print the pull progress
	decoder := json.NewDecoder(out)
	for {
		var pullStatus map[string]interface{}
		if err := decoder.Decode(&pullStatus); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		status, _ := pullStatus["status"].(string)
		if status != "" {
			fmt.Printf("%s: %s\n", imageName, status)
		}
	}
	return nil
}

// start a docker container
func DockerContainerCreate() {
	imageName := "govcmsextras/dnsmasq"

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		panic(err)
	}

	platform := platforms.Normalize(v1.Platform{
		Architecture: runtime.GOARCH,
		OS:           "linux",
	})

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, &platform, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}
