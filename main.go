package main

import (
	"context"
	"github.com/Zelayan/auto-deploy/docker"
	"github.com/docker/docker/client"
)

func main() {

	ctx := context.Background()
	myApi := docker.NewDockerApi(ctx, Client())
	myApi.CheckAndStartContainer(ctx, "test", "test:test")
}

func Client() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	return cli
}
