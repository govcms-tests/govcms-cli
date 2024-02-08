package docker

import (
	"context"
	"fmt"
	"io"
	"os"
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

// pull docker image
func DockerImagePull(imageName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		panic(err)
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
