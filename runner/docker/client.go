package docker

import (
	dockerclient "github.com/docker/docker/client"
	"github.com/w-h-a/pkg/telemetry/log"
)

type DockerClient interface {
}

type dockerClient struct {
	*dockerclient.Client
}

func NewDockerClient() DockerClient {
	client, err := dockerclient.NewClientWithOpts(dockerclient.WithHost(dockerclient.DefaultDockerHost))
	if err != nil {
		log.Fatal(err)
	}

	return &dockerClient{client}
}

type dockerMock struct{}

func NewDockerMock() DockerClient {
	return &dockerMock{}
}
