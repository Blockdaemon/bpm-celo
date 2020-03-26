package main

import (
	"github.com/Blockdaemon/bpm-sdk/pkg/docker"
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
	"github.com/blockdaemon/bpm-celo/configs"
	"github.com/blockdaemon/bpm-celo/pkg/tester"
)

var version string
var celoContainerName string = "celo"

const (
	celoContainerImage = "us.gcr.io/celo-testnet/celo-node:1.9"
	celoDataVolumeName = "celo"
	celoCmdFile        = "celo.dockercmd"

	collectorContainerName = "collector"
	collectorImage         = "docker.io/blockdaemon/celo-collector:BD-2901-get-rewards"
	collectorEnvFile       = "configs/collector.env"
)

func main() {

	description := "A Celo BPM Plugin"

	parameters := getParameters() // pass in subtype
	containers := getContainers() // pass in subtype

	templates := map[string]string{
		celoCmdFile:      configs.CeloCmdTpl,
		collectorEnvFile: configs.CollectorEnvTpl,
	}

	celoPlugin := plugin.NewDockerPlugin("celo", version, description, parameters, templates, containers)
	celoPlugin.Tester = tester.CeloTester{}

	plugin.Initialize(celoPlugin)
}

func getParameters() []plugin.Parameter {
	return []plugin.Parameter{
		{
			Name:        "subtype",
			Type:        plugin.ParameterTypeString,
			Description: "The type of node. Must be either `validator`, `proxy`, `fullnode`, `accounts` or `attestations`",
			Mandatory:   false,
			Default:     "fullnode",
		},
		{
			Name:        "networkid",
			Type:        plugin.ParameterTypeString,
			Description: "The current Celo network id",
			Mandatory:   true,
			Default:     "",
		},
		{
			Name:        "signer",
			Type:        plugin.ParameterTypeString,
			Description: "The signer address",
			Mandatory:   false,
			Default:     "",
		},
		{
			Name:        "bootnodes",
			Type:        plugin.ParameterTypeString,
			Description: "List of bootnodes to connect to",
			Mandatory:   false,
			Default:     "",
		},
		{
			Name:        "genesis",
			Type:        plugin.ParameterTypeString,
			Description: "Path to a genesis.json",
			Mandatory:   true,
		},
	}
}

func getContainers() []docker.Container {
	return []docker.Container{
		{
			Name:  "celoinit",
			Image: celoContainerImage,
			Cmd:   []string{"--nousb", "init", "/celo/genesis.json"},
			Mounts: []docker.Mount{
				{
					Type: "volume",
					From: celoDataVolumeName,
					To:   "/root/.celo",
				},
			},
			Restart:     "no",
			CollectLogs: false,
		},
		{
			Name:    celoContainerName,
			Image:   celoContainerImage,
			CmdFile: celoCmdFile,
			Mounts: []docker.Mount{
				{
					Type: "volume",
					From: celoDataVolumeName,
					To:   "/root/.celo",
				},
			},
			Ports: []docker.Port{
				{
					HostIP:        "0.0.0.0",
					HostPort:      "30333",
					ContainerPort: "30333",
					Protocol:      "tcp",
				},
				{
					HostIP:        "0.0.0.0",
					HostPort:      "30333",
					ContainerPort: "30333",
					Protocol:      "udp",
				},
				{
					HostIP:        "0.0.0.0",
					HostPort:      "30503",
					ContainerPort: "30503",
					Protocol:      "tcp",
				},
				{
					HostIP:        "0.0.0.0",
					HostPort:      "30503",
					ContainerPort: "30503",
					Protocol:      "udp",
				},
			},
			CollectLogs: true,
		},
		docker.Container{
			Name:        collectorContainerName,
			Image:       collectorImage,
			EnvFilename: collectorEnvFile,
			Mounts: []docker.Mount{
				{
					Type: "bind",
					From: "logs",
					To:   "/data/nodestate",
				},
			},
			CollectLogs: true,
		},
	}
}
