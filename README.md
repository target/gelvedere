# gelvedere

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)
[![release](https://img.shields.io/github/release/target/gelvedere.svg)](https://github.com/target/gelvedere/releases/latest)
[![watch](https://img.shields.io/github/watchers/target/gelvedere.svg?style=social)](https://github.com/target/gelvedere/watchers)
[![star](https://img.shields.io/github/stars/target/gelvedere.svg?style=social)](https://github.com/target/gelvedere/stargazers)

Cli to deploy a Jenkins master within the JAYS architecture

## How It Works

gelvedere provides a CLI for creating a Jenkins master within Docker swarm. Currently there are 2 types of input files required.

* `admin.json` - contains information specific to deploying the master within docker swarm
* `user.json` - contains information specific to Jenkins ACL configuration

## Getting Started

A sample command to run gelvedere:

```bash
gelvedere --user-config /jenkins/user-configs/test.json --admin-config /jenkins/admin-configs/test.json --domain acme.com
```

### Input files

A sample `user.json` file:

```json
{
  "name": "example",
  "admins": "target*Jenkins",
  "members": "",
  "team": "Jenkins"
}
```

A sample `admin.json` file:

```json
{
  "ghe_key": "1234",
  "ghe_secret": "56789",
  "port": "50000",
  "admin_ssh_pubkey": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDVr+0LAocyLbzzvQEdwjU8o+w0IYpR4R0uf2mswNYz6utcUVqHp5VXFog6YL4gYf0Q7naorLGh/zbROGHmBGAUngUbvy1vAnyiiBEjLPhW5k6iLy9f3N2lZyDQJ/VZYeRzfSeOPyEfd13MOjR8kB0zrodFa5j3fIToUrPmLytAVWplbF002jjJOTjwhFaknbdcVTzQ1LxhaOCaVjbEQyuFB3e8mB15kGEJOllnq4Un1HnG6wOcSx8IwP/E1JcmChfM3pPY2PIpYRqYaT4SYKGua+qke90aPNFl/k3j3J3yl2ZKGno/tJjj50sbTDgNz46uTLuLI2Eb6ETeI3d2Jy0Z jenkins@example.com",
  "size": "small",
  "image": "target/jenkins-docker-master:2.73.1-1"
}
```

Information on generating the values for the `user.json` and `admin.json` file can be found [here](https://github.com/target/jenkins-docker-master/tree/master/examples)

### Storage

gelvedere makes the assumption that you are using local storage for each master, but you can override the storage path with a environment variable or cli argument.

The below example will create a docker swarm service with a mount source of `/jenkins/stores` and a target of `/var/jenkins_home`.

```shell
gelvedere --user-config /jenkins/user-configs/stores.json --admin-config /jenkins/admin-configs/stores.json --mount-path /jenkins/stores
```

### Custom Environment Variables

By default gelvedere adds the following environment variables to the Docker swarm service:

* JENKINS_ACL_MEMBERS_admin
* JENKINS_ACL_MEMBERS_developer
* JENKINS_URL
* GHE_KEY
* GHE_SECRET
* ADMIN_SSH_PUBKEY
* JENKINS_SLAVE_AGENT_PORT
* JAVA_OPTS

More information on the above configuration options can be found [here](https://github.com/target/jenkins-docker-master/tree/master/examples#usage)

You can add additional variables by setting the following configuration in the `admin.json` input file.

```json
{
    "env_variables": {
        "<variable name>": "<variable value>",
        "<variable name>": "<variable value>"
    }
}
```

### Custom docker logging

By default gelvedere does not create Docker swarm services with [logging drivers](https://docs.docker.com/engine/admin/logging/overview/), but you can override that behavior by adding contents into the `admin.json` file.

The below example sends logs in the gelf format

```json
  "log_config": {
    "driver": "gelf",
    "gelf-address": "udp://gelf.example.com:12201",
    "tag": "custom-tag"
  }
```

## CLI Installation

### Install on Linux

```bash
curl -L https://github.com/target/gelvedere/releases/download/v0.1.2/gelvedere-linux-amd64.tgz | tar zx

sudo install -t /usr/local/bin gelvedere
```

### Install on macOS

```bash
curl -L https://github.com/target/gelvedere/releases/download/v0.1.2/gelvedere-darwin-amd64.tgz | tar zx

sudo cp gelvedere /usr/local/bin/
```
