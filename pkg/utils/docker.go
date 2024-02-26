package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/containerd/containerd/platforms"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"log"
	"runtime"
)

// DockerImageList list docker image(s)
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

// DockerImagePull Pulls a Docker image for a GovCMS service and prints the output in a nicer format.
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

// DockerContainerCreate starts a docker dockerContainer
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

func GetContainerIpByName(name string) string {
	containerID, err := GetContainerIdByName(name + "-govcms-1")
	if err != nil {
		log.Println(err)
	}
	ip, _ := GetContainerIP(containerID)
	return ip
}

// GetContainerIP returns the getContainerIP address of the dockerContainer.
func GetContainerIP(containerID string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	containerInspect, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}

	if len(containerInspect.NetworkSettings.Networks) > 0 {
		for _, port := range containerInspect.NetworkSettings.Ports {
			ip := port[0].HostIP + ":" + port[0].HostPort
			return ip, nil
		}
	}

	return "", fmt.Errorf("container %s has no network information", containerID)
}

func GetContainerIdByName(name string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return "", err
	}

	// Iterate through containers to find one with a matching name
	for _, cont := range containers {
		for _, cname := range cont.Names {
			if cname == "/"+name {
				return cont.ID, nil
			}
		}
	}
	return "", fmt.Errorf("container with name %s not found", name)

}

func CreateContainer() {

}
