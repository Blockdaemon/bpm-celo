package main

import (
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
	"github.com/blockdaemon/bpm-celo/pkg/celo"
	"github.com/blockdaemon/bpm-celo/pkg/tester"
)

const (
	celoContainerImage = "us.gcr.io/celo-testnet/celo-node:baklava"
	description        = "A Celo BPM Plugin"
	version            = "0.0.1"
)

func main() {

	c := celo.New(celoContainerImage)

	parameters := c.GetParameters()
	containers := c.GetContainers()
	templates := c.GetTemplates()

	celoPlugin := plugin.NewDockerPlugin("celo", version, description, parameters, templates, containers)
	celoPlugin.Tester = tester.CeloTester{}

	plugin.Initialize(celoPlugin)
}
