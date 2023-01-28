package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"os"
	"time"
)

type DockerAPI interface {
	List(ctx context.Context) ([]string, error)
	CreateContainer(ctx context.Context, name string) error
}

type Docker struct {
	ctx    context.Context
	client *client.Client
	n      int64
}

func NewDockerApi(ctx context.Context, client *client.Client) *Docker {
	return &Docker{
		ctx:    ctx,
		client: client,
	}
}

func (d *Docker) CreateContainer(ctx context.Context, name string) error {
	imageName := "test:test"
	_, err := d.client.ContainerCreate(ctx, &container.Config{
		Image:     imageName, //Docker基于该镜像创建容器
		Tty:       true,      //docker run 命令的-t
		OpenStdin: true,      //docker run命令的-i
		ExposedPorts: nat.PortSet{ //docker容器对外开放的端口
			"7070": struct{}{},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"7070": []nat.PortBinding{nat.PortBinding{ //docker容器映射到宿主机的端口
				HostIP:   "0.0.0.0",
				HostPort: "7070",
			}},
		},
	}, nil, nil, "test")

	if err != nil {
		return err
	}
	return nil

}

func (d *Docker) IsRunning(ctx context.Context, containerName string) bool {
	isRun := make(chan bool)
	var timer *time.Ticker
	go func(ctx context.Context) {
		for {
			//每n s检查一次容器是否运行

			timer = time.NewTicker(time.Duration(d.n) * time.Second)
			select {
			case <-timer.C:
				//获取正在运行的container list
				log.Printf("%s is checking the container[%s]is Runing??", os.Args[0], containerName)
				contTemp := d.getContainer(ctx, containerName, false)
				if contTemp.ID == "" {
					log.Print(":NO")
					//说明container没有运行
					isRun <- false
				} else {
					log.Print(":YES")
					//说明该container正在运行
					//go printConsole(ctx, cli, contTemp.ID)
					isRun <- true
				}
			}

		}
	}(ctx)
	return <-isRun
}

// 获取container
func (d *Docker) getContainer(ctx context.Context, containerName string, all bool) types.Container {
	containerList, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: all})
	if err != nil {
		panic(err)
	}
	var contTemp types.Container
	//找出名为“mygin-latest”的container并将其存入contTemp中
	for _, v1 := range containerList {
		for _, v2 := range v1.Names {
			if v2 == containerName {
				contTemp = v1
				break
			}
		}
	}
	return contTemp
}

// 启动容器
func (d *Docker) dstartContainer(ctx context.Context, containerID string, cli *client.Client) error {
	err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err == nil {
		log.Printf("success start container:%s\n", containerID)
	} else {
		log.Printf("failed to start container:%s!!!!!!!!!!!!!\n", containerID)
	}
	return err
}

func (d *Docker) List(ctx context.Context) ([]string, error) {
	var res []string
	list, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range list {
		res = append(res, list[i].Names...)
	}
	return res, nil
}
