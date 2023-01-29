package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
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
		n:      int64(5),
	}
}

func (d *Docker) CreateContainer(ctx context.Context, imageName string, containerName string) error {
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
	}, nil, nil, containerName)

	if err != nil {
		log.Printf("create container failed: %s, %s", containerName, err.Error())
		return err
	}
	return nil
}

func (d *Docker) IsRunning(ctx context.Context, containerName string) chan bool {
	isRun := make(chan bool)
	var timer *time.Ticker
	//每n s检查一次容器是否运行
	timer = time.NewTicker(time.Duration(d.n) * time.Second)
	go func(ctx context.Context) {
		for {
			select {
			case <-timer.C:
				//获取正在运行的container list
				log.Printf("%s is checking the container[%s]is Running??", os.Args[0], containerName)
				contTemp := d.GetContainer(ctx, containerName, false)
				if contTemp.ID == "" {
					log.Print(":NO")
					//说明container没有运行
					isRun <- false
				} else {
					log.Print(":YES")
					//说明该container正在运行
					go printConsole(ctx, d.client, contTemp.ID)
					isRun <- true
				}
			}

		}
	}(ctx)
	return isRun
}

// 获取container
func (d *Docker) GetContainer(ctx context.Context, containerName string, all bool) types.Container {
	containerList, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: all})
	if err != nil {
		panic(err)
	}
	var contTemp types.Container
	for _, v1 := range containerList {
		for _, v2 := range v1.Names {
			if v2 == "/"+containerName {
				contTemp = v1
				break
			}
		}
	}
	return contTemp
}

// 启动容器
func (d *Docker) startContainer(ctx context.Context, containerID string) error {
	err := d.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err == nil {
		log.Printf("success start container:%s\n", containerID)
	} else {
		log.Printf("failed to start container:%s!!!!!!!!!!!!!\n", containerID)
	}
	return err
}

func (d *Docker) CheckAndStartContainer(ctx context.Context, containerName string, imageName string) {
	for {
		select {
		case x := <-d.IsRunning(ctx, containerName):
			if !x {
				//该container没有在运行
				//获取所有的container查看该container是否存在
				contTemp := d.GetContainer(ctx, containerName, true)
				if contTemp.ID == "" {
					//该容器不存在，创建该容器
					log.Printf("the container name[%s] is not exists!!!!!!!!!!!!!\n", containerName)
					d.CreateContainer(ctx, imageName, containerName)
				} else {
					//该容器存在，启动该容器
					log.Printf("the container name[%s] is exists\n", containerName)
					d.startContainer(ctx, contTemp.ID)
				}
			}
		}
	}
}

// 将容器的标准输出输出到控制台中
func printConsole(ctx context.Context, cli *client.Client, id string) {
	//将容器的标准输出显示出来
	out, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	//容器内部的运行状态
	/*	status, err := cli.ContainerStats(ctx, id, true)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, status.Body)*/
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
