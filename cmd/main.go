package main

import (
<<<<<<< HEAD
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
	"github.com/blockdaemon/bpm-celo/pkg/celo"
	"github.com/blockdaemon/bpm-celo/pkg/tester"
)

const (
	celoContainerImage = "us.gcr.io/celo-testnet/celo-node:baklava"
	description        = "A Celo BPM Plugin"
	version            = "0.0.1"
=======
	"log"
	"os"

	"go.blockdaemon.com/bpm/celo/pkg/celo"
	"go.blockdaemon.com/bpm/celo/pkg/tester"
	"go.blockdaemon.com/bpm/sdk/pkg/plugin"
)

const (
	description = "A Celo BPM Plugin"
	version     = "0.0.2"
>>>>>>> release-0.0.2
)

func main() {

<<<<<<< HEAD
	c := celo.New(celoContainerImage)
=======
	c := celo.New()
>>>>>>> release-0.0.2

	parameters := c.GetParameters()
	containers := c.GetContainers()
	templates := c.GetTemplates()

	celoPlugin := plugin.NewDockerPlugin("celo", version, description, parameters, templates, containers)
	celoPlugin.Tester = tester.CeloTester{}

<<<<<<< HEAD
=======
	if c.Subtype != "attestation-service" {
		cmd := os.Args[1]
		if cmd == "start" {
			log.Println("Initialize genesis...")
			c.InitGenesis() // TODO handle erros, ffs (palmface)
		}
	}

>>>>>>> release-0.0.2
	plugin.Initialize(celoPlugin)
}
