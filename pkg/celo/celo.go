package celo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Blockdaemon/bpm-sdk/pkg/docker"
	"github.com/Blockdaemon/bpm-sdk/pkg/node"
	"github.com/Blockdaemon/bpm-sdk/pkg/plugin"
	"github.com/blockdaemon/bpm-celo/configs"
)

// Celo the main struct for this package
type Celo struct {
	image   string
	cmdFile string
	n       node.Node
	subtype string
}

// ICelo The Celo interface
type ICelo interface {
	GetParameters(string) []plugin.Parameter
	SetImage(string) bool
}

// New Returns a new Celo instance
func New(image string) *Celo {
	var c Celo

	n := getNode()

	c.cmdFile = "celo.dockercmd"
	c.image = image
	c.n = n
	c.subtype = n.StrParameters["subtype"]

	return &c
}

func getNode() node.Node {
	// load node.json
	var jsonfile string
	var n node.Node
	var err error

	for _, arg := range os.Args {
		if strings.Contains(arg, ".json") {
			jsonfile = arg
		}
	}
	if os.Args[1] != "meta" {
		n, err = node.Load(jsonfile)
		if err != nil {
			log.Fatalf("Unable to load node json: %s\n", err)
		}
	} else {
		n = node.New(jsonfile)
	}

	return n
}

// GetParameters returns parameters for a subtype, defaults to all
func (c *Celo) GetParameters() []plugin.Parameter {

	var params []plugin.Parameter
	subtype := c.subtype

	pSubtype := plugin.Parameter{
		Name:        "subtype",
		Type:        plugin.ParameterTypeString,
		Description: "The type of node. Must be either `validator`, `proxy`, `fullnode`, `accounts` or `attestations`",
		Mandatory:   false,
		Default:     "fullnode",
	}
	pNetworkID := plugin.Parameter{
		Name:        "networkid",
		Type:        plugin.ParameterTypeString,
		Description: "The current Celo network id",
		Mandatory:   false,
		Default:     "",
	}
	pSigner := plugin.Parameter{
		Name:        "signer",
		Type:        plugin.ParameterTypeString,
		Description: "The signer address",
		Mandatory:   false,
		Default:     "",
	}
	pBootnodes := plugin.Parameter{
		Name:        "bootnodes",
		Type:        plugin.ParameterTypeString,
		Description: "List of bootnodes to connect to",
		Mandatory:   false,
		Default:     "",
	}

	pKeystore := plugin.Parameter{
		Name:        "keystore-file",
		Type:        plugin.ParameterTypeString,
		Description: "Location of the signer keystore json",
		Mandatory:   false,
		Default:     "",
	}
	pKeypass := plugin.Parameter{
		Name:        "keystore-pass",
		Type:        plugin.ParameterTypeString,
		Description: "The password for the keystore json",
		Mandatory:   false,
		Default:     "",
	}
	pPort := plugin.Parameter{
		Name:        "port",
		Type:        plugin.ParameterTypeString,
		Description: "Port to listen to",
		Mandatory:   false,
		Default:     "30303",
	}
	pProxyInternal := plugin.Parameter{
		Name:        "proxy_internal",
		Type:        plugin.ParameterTypeString,
		Description: "The internal proxy ip, use external if none",
		Mandatory:   false,
		Default:     "",
	}
	pProxyExternal := plugin.Parameter{
		Name:        "proxy_external",
		Type:        plugin.ParameterTypeString,
		Description: "The external proxy ip",
		Mandatory:   false,
		Default:     "",
	}
	pEnode := plugin.Parameter{
		Name:        "enode",
		Type:        plugin.ParameterTypeString,
		Description: "The proxy enode id",
		Mandatory:   false,
		Default:     "",
	}
	pRpcaddr := plugin.Parameter{
		Name:        "rpcaddr",
		Type:        plugin.ParameterTypeString,
		Description: "The rpcaddr ip address",
		Mandatory:   false,
		Default:     "0.0.0.0",
	}
	pLightServe := plugin.Parameter{
		Name:        "light_serve",
		Type:        plugin.ParameterTypeString,
		Description: "light.serve",
		Mandatory:   false,
		Default:     "10",
	}
	pMaxpeers := plugin.Parameter{
		Name:        "maxpeers",
		Type:        plugin.ParameterTypeString,
		Description: "The max peers `light.maxpeers`",
		Mandatory:   false,
		Default:     "10",
	}
	pAccount := plugin.Parameter{
		Name:        "account",
		Type:        plugin.ParameterTypeString,
		Description: "The account to send rewards to",
		Mandatory:   false,
		Default:     "",
	}

	switch subtype {

	case "proxy":
		pNetworkID.Mandatory = true
		pSigner.Mandatory = true
		pBootnodes.Mandatory = true
		params = []plugin.Parameter{
			pSubtype,
			pNetworkID,
			pSigner,
			pBootnodes,
		}
	case "validator":
		pNetworkID.Mandatory = true
		pSigner.Mandatory = true
		pKeystore.Mandatory = true
		pKeypass.Mandatory = true
		pProxyInternal.Mandatory = true
		pProxyExternal.Mandatory = true
		pEnode.Mandatory = true
		params = []plugin.Parameter{
			pSubtype,
			pNetworkID,
			pSigner,
			pKeystore,
			pKeypass,
			pPort,
			pProxyInternal,
			pProxyExternal,
			pEnode,
		}
	case "fullnode":
		params = []plugin.Parameter{
			pSubtype,
			pNetworkID,
			pPort,
			pBootnodes,
			pRpcaddr,
			pLightServe,
			pMaxpeers,
			pAccount,
			pPort,
		}
	default: // show all params
		params = []plugin.Parameter{
			pSubtype,
			pNetworkID,
			pSigner,
			pKeystore,
			pKeypass,
			pPort,
			pProxyInternal,
			pProxyExternal,
			pEnode,
			pBootnodes,
			pRpcaddr,
			pLightServe,
			pMaxpeers,
			pAccount,
		}
	}

	return params
}

// GetContainers returns containers for a subtype, defaults to all
func (c *Celo) GetContainers() []docker.Container {

	collectorContainerName := "collector"
	collectorImage := "docker.io/blockdaemon/celo-collector:0.0.3"
	collectorEnvFile := "configs/collector.env"

	var containers []docker.Container
	subtype := c.subtype
	n := c.n
	volumeName := "celo-" + subtype

	cBootnodes := docker.Container{
		Name:    "celoinit",
		Image:   c.image,
		Cmd:     []string{"--nousb", "init", "/celo/genesis.json"},
		Restart: "no",
		Mounts: []docker.Mount{
			{
				Type: "volume",
				From: volumeName,
				To:   "/root/.celo",
			},
		},
		CollectLogs: false,
	}
	cProxy := docker.Container{
		Name:    "proxy",
		Image:   c.image,
		CmdFile: c.cmdFile,
		Mounts: []docker.Mount{
			{
				Type: "volume",
				From: volumeName,
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
	}
	cValidator := docker.Container{
		Name:    "validator",
		Image:   c.image,
		CmdFile: c.cmdFile,
		Mounts: []docker.Mount{
			{
				Type: "volume",
				From: volumeName,
				To:   "/root/.celo",
			},
			{
				Type: "bind",
				From: "./configs",
				To:   "/root/.celo/configs",
			},
		},
		Ports: []docker.Port{
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: n.StrParameters["port"],
				Protocol:      "tcp",
			},
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: n.StrParameters["port"],
				Protocol:      "udp",
			},
		},
		CollectLogs: true,
	}
	cFullnode := docker.Container{
		Name:    "fullnode",
		Image:   c.image,
		CmdFile: c.cmdFile,
		Mounts: []docker.Mount{
			{
				Type: "volume",
				From: volumeName,
				To:   "/root/.celo",
			},
		},
		Ports: []docker.Port{
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: n.StrParameters["port"],
				Protocol:      "tcp",
			},
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: n.StrParameters["port"],
				Protocol:      "udp",
			},
		},
		CollectLogs: true,
	}
	cNodestate := docker.Container{
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
	}

	switch subtype {
	case "proxy":
		containers = []docker.Container{
			cBootnodes,
			cProxy,
			cNodestate,
		}
	case "validator":
		containers = []docker.Container{
			cBootnodes,
			cValidator,
			cNodestate,
		}
	case "fullnode":
		containers = []docker.Container{
			cBootnodes,
			cFullnode,
			cNodestate,
		}
	default:
		containers = []docker.Container{
			cBootnodes,
			cProxy,
		}
	}

	return containers
}

// GetTemplates Returns the templates for current node
func (c *Celo) GetTemplates() map[string]string {
	subtype := c.subtype
	dockerCmd := c.getDockerCmd()

	templates := map[string]string{
		"celo.dockercmd":        dockerCmd,
		"configs/collector.env": configs.CollectorEnvTpl,
	}

	if subtype == "validator" {
		ks := c.getKeystore()
		templates["configs/keystore/"+ks.filename] = ks.json // string
		templates["configs/password.secret"] = ks.pass
	}

	return templates
}

func (c *Celo) getDockerCmd() string {
	subtype := c.subtype

	dockerCmd := ``
	switch subtype {
	case "proxy":
		dockerCmd = configs.ProxyCmdTpl
	case "validator":
		dockerCmd = configs.ValidatorCmdTpl
	case "fullnode":
		dockerCmd = configs.FullnodeCmdTpl
	default:
		dockerCmd = "help" // docker command required by sdk?
	}

	return dockerCmd
}

type keystore struct {
	filename string
	pass     string
	json     string
}

func (c *Celo) getKeystore() keystore {

	n := c.n
	file := n.StrParameters["keystore-file"]
	pass := n.StrParameters["keystore-pass"]

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error opening keystore file: %s: %s\n", file, err)
	}
	p, err := ioutil.ReadFile(pass)
	if err != nil {
		log.Printf("Error opening password file: %s: %s\n", pass, err)
	}

	ks := keystore{
		filename: filepath.Base(file),
		pass:     string(p),
		json:     string(content),
	}

	nodeDir := n.NodeDirectory()
	targetDir := nodeDir + "/configs/keystore"
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating keystore directory: %s\n", err)
	}

	return ks
}
