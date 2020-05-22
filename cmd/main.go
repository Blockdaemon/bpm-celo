package main

import (
	"log"
	"os"

	"github.com/blockdaemon/bpm-celo/pkg/celo"
	"github.com/blockdaemon/bpm-celo/pkg/tester"
	"go.blockdaemon.com/bpm/sdk/pkg/plugin"
)

const (
	description = "A Celo BPM Plugin"
	version     = "0.0.1"
)

func main() {

	c := celo.New()

	parameters := c.GetParameters()
	containers := c.GetContainers()
	templates := c.GetTemplates()

	celoPlugin := plugin.NewDockerPlugin("celo", version, description, parameters, templates, containers)
	celoPlugin.Tester = tester.CeloTester{}

	if c.Subtype != "attestation-service" {
		cmd := os.Args[1]
		if cmd == "start" {
			log.Println("Initialize genesis...")
			c.InitGenesis() // TODO handle erros, ffs (palmface)
		}
	}

	plugin.Initialize(celoPlugin)
}
