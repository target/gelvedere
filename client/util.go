package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/target/gelvedere/model"

	"github.com/docker/docker/api/types/swarm"
)

// CheckName checks if the docker swarm service name is already in use
func CheckName(serviceName string, services []swarm.Service) error {
	for _, service := range services {
		if service.Spec.Name == serviceName {
			return fmt.Errorf("service %v is already in use", serviceName)
		}
	}

	return nil
}

// CheckPort checks if the port is already in use
func CheckPort(servicePort int, services []swarm.Service) error {
	for _, service := range services {
		if service.Endpoint.Ports != nil {
			for _, v := range service.Endpoint.Ports {
				if v.PublishedPort == uint32(servicePort) {
					return fmt.Errorf("port %v is already in use", servicePort)
				}
			}
		}
	}
	return nil
}

// GetAdminJSON returns a json config for admins
func GetAdminJSON(file string) (*model.AdminConfig, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON contents to AdminConfig
	var ac model.AdminConfig
	err = json.Unmarshal(bytes, &ac)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}

// GetUserJSON returns a json config for users
func GetUserJSON(file string) (*model.UserConfig, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON contents to UserConfig
	var uc model.UserConfig
	err = json.Unmarshal(bytes, &uc)
	if err != nil {
		return nil, err
	}

	return &uc, nil
}
