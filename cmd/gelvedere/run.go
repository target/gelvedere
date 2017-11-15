package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/target/gelvedere/client"
	"github.com/target/gelvedere/model"
	cli "gopkg.in/urfave/cli.v1"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

func run(c *cli.Context) error {
	// debug level if requested by user
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	ac, err := client.GetAdminJSON(c.String("admin-config"))
	if err != nil {
		return err
	}

	uc, err := client.GetUserJSON(c.String("user-config"))
	if err != nil {
		return err
	}

	mountPath := c.String("mount-path")
	if mountPath == "" {
		mountPath = fmt.Sprintf("/jenkins/qtr_%s", uc.Name)

		if _, err = os.Stat(mountPath); os.IsNotExist(err) {
			return fmt.Errorf("Source directory %s does not exist", mountPath)
		}
	}

	err = Validate(uc, ac)
	if err != nil {
		return err
	}

	serviceSpec, err := GetServiceSpec(uc, ac, c)
	if err != nil {
		return err
	}

	resp, err := client.CreateDockerSwarmService(serviceSpec)
	if err != nil {
		return err
	}

	log.Info(resp)

	return nil
}

// GetServiceSpec gets the configuration to create a docker swarm service
func GetServiceSpec(uc *model.UserConfig, ac *model.AdminConfig, c *cli.Context) (swarm.ServiceSpec, error) {
	duration := time.Duration(10) * time.Second

	jenkinsURL := c.String("url")
	if jenkinsURL == "" {
		if ac.Region == "" {
			jenkinsURL = fmt.Sprintf("https://%s.%s.%s", uc.Name, c.String("subdomain"), c.String("domain"))
		} else {
			jenkinsURL = fmt.Sprintf("https://%s.%s.%s.%s", uc.Name, ac.Region, c.String("subdomain"), c.String("domain"))
		}
	}

	mountPath := c.String("mount-path")
	if mountPath == "" {
		mountPath = fmt.Sprintf("/jenkins/qtr_%s", uc.Name)
	}

	resources := getResources(ac.Size)
	// convert bytes to format for xms and xmx
	mem := bytefmt.ByteSize(uint64(resources.Reservations.MemoryBytes))
	// stip G so we can convert mem to int
	mem = mem[:len(mem)-1]
	// convert string to int so we can do math
	xm, err := strconv.Atoi(mem)
	// reduce the memory by 1 since we want xmx to be 1g less than the system max
	xm = xm - 1
	if err != nil {
		return swarm.ServiceSpec{}, err
	}
	memOpts := fmt.Sprintf("-Xms%vg -Xmx%vg", xm, xm)

	port, _ := strconv.Atoi(ac.Port)

	containerEnv := buildContainerEnv(uc, ac, jenkinsURL, memOpts)

	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: uc.Name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: ac.Image,
				Env:   containerEnv,
				Mounts: []mount.Mount{
					mount.Mount{
						Source: mountPath,
						Target: "/var/jenkins_home",
					},
				},
			},
			Resources: resources,
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyCondition("on-failure"),
				Delay:     &duration,
			},
			Placement: &swarm.Placement{
				Constraints: []string{"node.role != manager"},
			},
		},
		Networks: []swarm.NetworkAttachmentConfig{
			swarm.NetworkAttachmentConfig{Target: "jenkins"},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Mode: swarm.ResolutionModeVIP,
			Ports: []swarm.PortConfig{
				swarm.PortConfig{
					Protocol:      swarm.PortConfigProtocolTCP,
					TargetPort:    uint32(port),
					PublishedPort: uint32(port),
					PublishMode:   swarm.PortConfigPublishModeIngress,
				},
			},
		},
	}

	if ac.LogConfig["driver"] != "" {
		logconfig := make(map[string]string)
		for k, v := range ac.LogConfig {
			if k != "driver" {
				logconfig[k] = v
			}
		}
		if logconfig["tag"] == "" {
			if ac.Region == "" {
				logconfig["tag"] = fmt.Sprintf("jenkins-%s", uc.Name)
			} else {
				logconfig["tag"] = fmt.Sprintf("jenkins-%s-%s", ac.Region, uc.Name)
			}
		}
		serviceSpec.TaskTemplate.LogDriver = &swarm.Driver{
			Name:    ac.LogConfig["driver"],
			Options: logconfig,
		}
	}

	return serviceSpec, nil
}

func getResources(size string) *swarm.ResourceRequirements {
	var resources *swarm.ResourceRequirements

	switch size {
	case "medium":
		resources = &swarm.ResourceRequirements{
			// 1 cpu and 9G ram
			Reservations: &swarm.Resources{NanoCPUs: 1000000000, MemoryBytes: 9663676416},
			// 4 cpu and 20G ram
			Limits: &swarm.Resources{NanoCPUs: 4000000000, MemoryBytes: 21474836480},
		}
	case "large":
		resources = &swarm.ResourceRequirements{
			// 2 cpu and 9G ram
			Reservations: &swarm.Resources{NanoCPUs: 2000000000, MemoryBytes: 9663676416},
			// 8 cpu and 20G ram
			Limits: &swarm.Resources{NanoCPUs: 8000000000, MemoryBytes: 21474836480},
		}
	case "small":
		fallthrough
	default:
		resources = &swarm.ResourceRequirements{
			// .5 cpu and 5G ram
			Reservations: &swarm.Resources{NanoCPUs: 500000000, MemoryBytes: 5368709120},
			// 4 cpu and 20G ram
			Limits: &swarm.Resources{NanoCPUs: 4000000000, MemoryBytes: 21474836480},
		}
	}

	return resources
}

func buildContainerEnv(uc *model.UserConfig, ac *model.AdminConfig, jenkinsURL, memoryOpts string) []string {
	var envs []string

	envs = append(envs, "JENKINS_ACL_MEMBERS_admin"+"="+uc.Admins)
	envs = append(envs, "JENKINS_ACL_MEMBERS_developer"+"="+uc.Members)
	envs = append(envs, "JENKINS_URL"+"="+jenkinsURL)
	envs = append(envs, "GHE_KEY"+"="+ac.GheKey)
	envs = append(envs, "GHE_SECRET"+"="+ac.GheSecret)
	envs = append(envs, "ADMIN_SSH_PUBKEY"+"="+ac.AdminSSH)
	envs = append(envs, "JENKINS_SLAVE_AGENT_PORT"+"="+ac.Port)
	envs = append(envs, "JAVA_OPTS"+"="+memoryOpts)

	for k, v := range ac.EnvVars {
		envs = append(envs, k+"="+v)
	}

	return envs
}
