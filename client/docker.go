package client

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// GetDockerSwarmServices gets docker swarm services
func GetDockerSwarmServices() ([]swarm.Service, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return nil, err
	}

	return services, nil
}

// CreateDockerSwarmService creates docker swarm service
func CreateDockerSwarmService(serviceSpec swarm.ServiceSpec) (*types.ServiceCreateResponse, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	resp, err := cli.ServiceCreate(ctx, serviceSpec, types.ServiceCreateOptions{})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
