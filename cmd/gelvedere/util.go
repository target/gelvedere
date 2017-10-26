package main

import (
	"strconv"

	"github.com/target/gelvedere/client"
	"github.com/target/gelvedere/model"
)

// Validate calls the necessary functions to validate the config files before
// creating the swarm service spec
func Validate(uc *model.UserConfig, ac *model.AdminConfig) error {
	services, err := client.GetDockerSwarmServices()
	if err != nil {
		return err
	}

	err = client.CheckName(uc.Name, services)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(ac.Port)
	if err != nil {
		return err
	}

	err = client.CheckPort(port, services)
	if err != nil {
		return err
	}

	return nil
}
