package celo

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.blockdaemon.com/bpm/celo/configs"
	"go.blockdaemon.com/bpm/sdk/pkg/docker"
	"go.blockdaemon.com/bpm/sdk/pkg/node"
	"go.blockdaemon.com/bpm/sdk/pkg/plugin"
)

// Celo the main struct for this package
type Celo struct {
	image            string
	imageAttestation string
	networkID        string
	cmdFile          string
	n                node.Node
	Subtype          string
}

// ICelo The Celo interface
type ICelo interface {
	GetParameters(string) []plugin.Parameter
	SetImage(string) bool
}

// New Returns a new Celo instance
func New() *Celo {
	var c Celo

	n := buildNode()

	// get the images
	if n.StrParameters["network"] == "baklava" {
		c.image = "us.gcr.io/celo-testnet/celo-node:baklava"
		c.imageAttestation = "us.gcr.io/celo-testnet/celo-monorepo:attestation-service-baklava"
		c.networkID = "40120"
	} else if n.StrParameters["network"] == "mainnet" {
		c.image = "us.gcr.io/celo-testnet/celo-node:rc1"
		c.imageAttestation = "us.gcr.io/celo-testnet/celo-monorepo:attestation-service-rc1"
		c.networkID = "42220"
	}

	c.cmdFile = "celo.dockercmd"
	c.n = n
	c.Subtype = n.StrParameters["subtype"]

	return &c
}

// GetParameters returns parameters for a subtype, defaults to all
func (c *Celo) GetParameters() []plugin.Parameter {

	var params []plugin.Parameter
	subtype := c.Subtype

	pSubtype := plugin.Parameter{
		Name:        "subtype",
		Type:        plugin.ParameterTypeString,
		Description: "The type of node. Must be either `validator`, `proxy` or `fullnode`",
		Mandatory:   false,
		Default:     "fullnode",
	}
	pNetwork := plugin.Parameter{
		Name:        "network",
		Type:        plugin.ParameterTypeString,
		Description: "Mainnet or baklava testnet",
		Mandatory:   true,
		Default:     "baklava",
	}
	pNetworkID := plugin.Parameter{
		Name:        "networkid",
		Type:        plugin.ParameterTypeString,
		Description: "The current Celo network id",
		Mandatory:   true,
		Default:     c.networkID,
	}
	pSigner := plugin.Parameter{
		Name:        "signer",
		Type:        plugin.ParameterTypeString,
		Description: "The signer address",
		Mandatory:   false,
		Default:     "",
	}
	pNoUSB := plugin.Parameter{
		Name:        "nousb",
		Type:        plugin.ParameterTypeString,
		Description: "Boolean. Wether to expect usb connections, eg ledger",
		Mandatory:   false,
		Default:     "true",
	}
	pValidator := plugin.Parameter{
		Name:        "validator",
		Type:        plugin.ParameterTypeString,
		Description: "The validator address",
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
	pRPCPort := plugin.Parameter{
		Name:        "rpcport",
		Type:        plugin.ParameterTypeString,
		Description: "The rpc port for the hostt",
		Mandatory:   false,
		Default:     "8545",
	}
	pLightServe := plugin.Parameter{
		Name:        "light_serve",
		Type:        plugin.ParameterTypeString,
		Description: "light.serve",
		Mandatory:   false,
		Default:     "90",
	}
	pLightMaxpeers := plugin.Parameter{
		Name:        "light_maxpeers",
		Type:        plugin.ParameterTypeString,
		Description: "The max peers `light.maxpeers`",
		Mandatory:   false,
		Default:     "1000",
	}
	pMaxpeers := plugin.Parameter{
		Name:        "maxpeers",
		Type:        plugin.ParameterTypeString,
		Description: "The max peers `light.maxpeers`",
		Mandatory:   false,
		Default:     "1100",
	}
	pAccount := plugin.Parameter{
		Name:        "account",
		Type:        plugin.ParameterTypeString,
		Description: "The account to send rewards to",
		Mandatory:   false,
		Default:     "",
	}
	// pCeloCommands := plugin.Parameter{
	// 	Name:        "celo",
	// 	Type:        plugin.ParameterTypeString,
	// 	Description: "Extra commands for container. Example: `--celo=\"--rpcapi web,personal,debug --rpcport 1234 --rpchost 0.0.0.0\"",
	// 	Mandatory:   false,
	// 	Default:     "",
	// }
	pDBHost := plugin.Parameter{
		Name:        "database",
		Type:        plugin.ParameterTypeString,
		Description: "Database URL for attestation service",
		Mandatory:   false,
		Default:     "",
	}
	pDBPassword := plugin.Parameter{
		Name:        "db_password",
		Type:        plugin.ParameterTypeString,
		Description: "Database password for attestation service postgres",
		Mandatory:   false,
		Default:     "",
	}
	pDBUser := plugin.Parameter{
		Name:        "db_user",
		Type:        plugin.ParameterTypeString,
		Description: "Database user for attestation service postgres",
		Mandatory:   false,
		Default:     "",
	}
	pAttNode := plugin.Parameter{
		Name:        "node_url",
		Type:        plugin.ParameterTypeString,
		Description: "Attestation node url, eg http://bpm-flower-pot-1234-attestattion-node:8545",
		Mandatory:   false,
		Default:     "",
	}
	pTwilioServiceSID := plugin.Parameter{
		Name:        "twilio_service_sid",
		Type:        plugin.ParameterTypeString,
		Description: "Twilio messaging service SID for attestation services",
		Mandatory:   false,
		Default:     "",
	}
	pTwilioAccountSID := plugin.Parameter{
		Name:        "twilio_account_sid",
		Type:        plugin.ParameterTypeString,
		Description: "Twilio account SID for attesation service",
		Mandatory:   false,
		Default:     "",
	}
	pTwilioBlacklist := plugin.Parameter{
		Name:        "twilio_blacklist",
		Type:        plugin.ParameterTypeString,
		Description: "Twilio blacklist for attesation service",
		Mandatory:   false,
		Default:     "",
	}
	pTwilioAuthToken := plugin.Parameter{
		Name:        "twilio_auth_token",
		Type:        plugin.ParameterTypeString,
		Description: "Auth token for Twilio",
		Mandatory:   false,
		Default:     "",
	}

	switch subtype {

	case "proxy":
		pSigner.Mandatory = true
		pBootnodes.Mandatory = true
		params = []plugin.Parameter{
			pNetwork,
			pSubtype,
			pNetworkID,
			pSigner,
			pBootnodes,
			// pCeloCommands,
		}
	case "validator":
		pSigner.Mandatory = true
		pKeystore.Mandatory = true
		pKeypass.Mandatory = true
		pProxyInternal.Mandatory = true
		pProxyExternal.Mandatory = true
		pEnode.Mandatory = true
		params = []plugin.Parameter{
			pNetwork,
			pSubtype,
			pNetworkID,
			pSigner,
			pKeystore,
			pKeypass,
			pPort,
			pProxyInternal,
			pProxyExternal,
			pEnode,
			// pCeloCommands,
		}
	case "fullnode":
		pBootnodes.Mandatory = true
		pAccount.Mandatory = true
		params = []plugin.Parameter{
			pNetwork,
			pSubtype,
			pNetworkID,
			pPort,
			pBootnodes,
			pRpcaddr,
			pRPCPort,
			pLightServe,
			pLightMaxpeers,
			pMaxpeers,
			pAccount,
			pPort,
			// pCeloCommands,
			pNoUSB,
		}
	case "attestation-node":
		pSigner.Mandatory = true
		pKeystore.Mandatory = true
		pKeypass.Mandatory = true
		pBootnodes.Mandatory = true
		params = []plugin.Parameter{
			pNetwork,
			pNetworkID,
			pSigner,
			pKeystore,
			pKeypass,
			pBootnodes,
			pRpcaddr,
			// pCeloCommands,
		}
	case "attestation-service":
		pSigner.Mandatory = true
		pValidator.Mandatory = true
		pAttNode.Mandatory = true
		pDBUser.Mandatory = true
		pDBPassword.Mandatory = true
		pTwilioServiceSID.Mandatory = true
		pTwilioAccountSID.Mandatory = true
		pTwilioAuthToken.Mandatory = true

		// get default db url
		if c.n.StrParameters["db_host"] == "" {
			c.getDBUrl()
		}

		params = []plugin.Parameter{
			pNetwork,
			pSigner,
			pValidator,
			pDBHost,
			pDBPassword,
			pAttNode,
			pTwilioServiceSID,
			pTwilioAccountSID,
			pTwilioAuthToken,
			pPort,
			pTwilioBlacklist,
			// pCeloCommands,
		}

	default: // show all params so they appear in the bpm manifest
		params = []plugin.Parameter{
			pSubtype,
			pNetwork,
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
			pRPCPort,
			pNoUSB,
			pLightServe,
			pLightMaxpeers,
			pMaxpeers,
			pAccount,
			// pCeloCommands,
			pValidator,
			pDBHost,
			pDBPassword,
			pDBUser,
			pAttNode,
			pTwilioServiceSID,
			pTwilioAccountSID,
			pTwilioAuthToken,
			pTwilioBlacklist,
		}
	}

	// // explode celo commands into lines
	// if c.n.StrParameters["celo"] != "" {
	// 	fmt.Println("regexp triggered")
	// 	reg := regexp.MustCompile(" ")
	// 	c.n.StrParameters["celo"] = reg.ReplaceAllString(c.n.StrParameters["celo"], "\n")
	// }
	// fmt.Printf("celocmd: %v\n", c.n.StrParameters["celo"])

	return params
}

// GetContainers returns containers for a subtype, defaults to all
func (c *Celo) GetContainers() []docker.Container {

	collectorContainerName := "collector"
	collectorImage := "docker.io/blockdaemon/celo-collector:0.0.5"
	collectorEnvFile := "configs/collector.env"
	postgresEnvFile := "configs/postgres.env"

	var containers []docker.Container
	subtype := c.Subtype
	datadir := c.n.StrParameters["data-dir"]
	n := c.n

	cProxy := docker.Container{
		Name:    "proxy",
		Image:   c.image,
		CmdFile: c.cmdFile,
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: datadir,
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
				Type: "bind",
				From: datadir,
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
				Type: "bind",
				From: datadir,
				To:   "/root/.celo",
			},
		},
		Ports: []docker.Port{
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: "30303",
				Protocol:      "tcp",
			},
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: "30303",
				Protocol:      "udp",
			},
			{
				HostIP:        c.n.StrParameters["rpcaddr"],
				HostPort:      c.n.StrParameters["rpcport"],
				ContainerPort: "8545",
			},
		},
		CollectLogs: true,
	}

	cAttestation := docker.Container{
		Name:    "attestation-node",
		Image:   c.image,
		CmdFile: c.cmdFile,
		Ports: []docker.Port{
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["rpcport"],
				ContainerPort: "8545",
				Protocol:      "tcp",
			},
		},
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: datadir,
				To:   "/root/.celo",
			},
			{
				Type: "bind",
				From: "./configs",
				To:   "/root/.celo/configs",
			},
		},
		CollectLogs: true,
	}
	cAttestationService := docker.Container{
		Name:    "attestation-service",
		Image:   c.imageAttestation,
		CmdFile: c.cmdFile,
		Ports: []docker.Port{
			{
				HostIP:        "0.0.0.0",
				HostPort:      n.StrParameters["port"],
				ContainerPort: n.StrParameters["port"],
				Protocol:      "tcp",
			},
		},
		CollectLogs: false,
		EnvFilename: "configs/attestation-service.env",
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
	cPostgres := docker.Container{
		Name:        "attestation-postgres",
		Image:       "docker.io/library/postgres:13",
		EnvFilename: postgresEnvFile,
		CollectLogs: false,
	}

	switch subtype {
	case "proxy":
		containers = []docker.Container{
			cProxy,
			cNodestate,
		}
	case "validator":
		containers = []docker.Container{
			cValidator,
			cNodestate,
		}
	case "fullnode":
		containers = []docker.Container{
			cFullnode,
			cNodestate,
		}
	case "attestation-node":
		containers = []docker.Container{
			cAttestation,
		}
	case "attestation-service":
		containers = []docker.Container{
			cPostgres,
			cAttestationService,
		}
	default:
		containers = []docker.Container{
			cProxy,
		}
	}

	return containers
}

func (c *Celo) GetNode() node.Node {
	return c.n
}

// GetTemplates Returns the templates for current node
func (c *Celo) GetTemplates() map[string]string {
	subtype := c.Subtype
	dockerCmd := c.getDockerCmd()

	templates := map[string]string{
		"celo.dockercmd": dockerCmd,
	}

	if subtype != "attestation-service" {
		templates["configs/collector.env"] = configs.CollectorEnvTpl
	}
	if subtype == "attestation-service" {

		if c.n.StrParameters["db_host"] == "" {
			c.getDBUrl()
		}
		dbURL := c.n.StrParameters["db_host"]

		templates["configs/attestation-service.env"] = strings.Replace(configs.AttesetationServiceEnvs, "{{ .Node.StrParameters.db_host }}", dbURL, -1)
		templates["configs/postgres.env"] = configs.PostgresEnvs
	}
	if subtype == "validator" || subtype == "attestation-node" {
		ks := c.getKeystore()
		templates["configs/keystore/"+ks.filename] = ks.json // string
		templates["configs/.password.secret"] = ks.pass
	}

	return templates
}

// InitGenesis Call `geth init /celo/genesis.json` in mounted dir to provision a Celo node.
func (c *Celo) InitGenesis() (bool, error) {

	bm, err := docker.NewBasicManager(c.n)
	if err != nil {
		return false, err
	}

	container := docker.Container{
		Name:  "celoinit",
		Image: c.image,
		Cmd:   []string{"--nousb", "init", "/celo/genesis.json"},
		Mounts: []docker.Mount{
			{
				Type: "bind",
				From: c.n.StrParameters["data-dir"],
				To:   "/root/.celo",
			},
		},
		CollectLogs: false,
	}

	ctx := context.Background()
	stdOut, err := bm.RunTransientContainer(ctx, container)
	if err != nil {
		return false, err
	}

	reg := regexp.MustCompile(`(Successfully\swrote\sgenesis\sstate)`)
	status := reg.FindSubmatch([]byte(stdOut))
	if len(status[1]) <= 0 {
		return false, fmt.Errorf("Unkown Error")
	}

	return true, nil
}

func (c *Celo) getDockerCmd() string {
	subtype := c.Subtype

	dockerCmd := ``
	switch subtype {
	case "proxy":
		dockerCmd = configs.ProxyCmdTpl
	case "validator":
		dockerCmd = configs.ValidatorCmdTpl
	case "fullnode":
		dockerCmd = configs.FullnodeCmdTpl
	case "attestation-node":
		dockerCmd = configs.AttestationCmdTpl
	case "attestation-service":
		dockerCmd = configs.AttestationServiceCmdTpl
	default:
		dockerCmd = "--help" // docker command required by sdk?
	}

	return strings.Replace(dockerCmd, "{{ .Node.StrParameters.networkid }}", c.networkID, -1)
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

func (c *Celo) getDBUrl() {
	postgres := "bpm-" + c.n.ID + "-attestation-postgres" + ":5432"
	c.n.StrParameters["db_host"] = "postgres://" + c.n.StrParameters["db_user"] + ":" + c.n.StrParameters["db_password"] + "@" + postgres
}

func buildNode() node.Node {
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
