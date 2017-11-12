// go build -o release/gelvedere ./cmd/gelvedere
package main

import (
	"fmt"
	"os"

	"github.com/target/gelvedere/version"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Version = version.Version.String()
	app.Name = "gelvedere"
	app.Action = run

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug level logging",
			EnvVar: "JENKINS_DEBUG",
		},
		cli.StringFlag{
			Name:   "user-config",
			Usage:  "path to json user configuration file",
			EnvVar: "JENKINS_USER_CONFIG",
			Value:  "user.json",
		},
		cli.StringFlag{
			Name:   "admin-config",
			Usage:  "path to json admin configuration file",
			EnvVar: "JENKINS_ADMIN_CONFIG",
			Value:  "admin.json",
		},
		cli.StringFlag{
			Name:   "mount-path",
			Usage:  "source mount path on host",
			EnvVar: "JENKINS_MOUNT_PATH",
		},
		cli.StringFlag{
			Name:   "url",
			Usage:  "url to reach jenkins master",
			EnvVar: "JENKINS_URL",
		},
		cli.StringFlag{
			Name:   "domain",
			Usage:  "jenkins host domain",
			EnvVar: "JENKINS_DOMAIN",
			Value:  "example.com",
		},
		cli.StringFlag{
			Name:   "subdomain",
			Usage:  "jenkins host subdomain",
			EnvVar: "JENKINS_SUBDOMAIN",
			Value:  "jenkins",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
