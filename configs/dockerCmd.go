package configs

const (
	// CollectorEnvTpl tempalte for collector command
	CollectorEnvTpl = `SERVICE_PORT=9933
SERVICE_HOST={{ .Node.NamePrefix }}proxy`

	// ProxyCmdTpl the celo command for running proxies
	ProxyCmdTpl = `--verbosity=3
--networkid={{ .Node.StrParameters.networkid }}
--syncmode=full
--proxy.proxy
--proxy.proxiedvalidatoraddress={{ .Node.StrParameters.signer }}
--proxy.internalendpoint=:30503
--etherbase={{ .Node.StrParameters.signer }}
--bootnodes={{ .Node.StrParameters.bootnodes }}
--ethstats={{ .Node.ID }}-proxy@baklava-ethstats.celo-testnet.org
`

	// ValidatorCmdTpl the celo command for running validator
	ValidatorCmdTpl = `--verbosity=3
--networkid={{ .Node.StrParameters.networkid }}
--syncmode=full
--mine
--port={{ .Node.StrParameters.port }}
--istanbul.blockperiod=5
--istanbul.requesttimeout=3000
--etherbase={{ .Node.StrParameters.signer }}
--nodiscover
--proxy.proxied
--proxy.proxyenodeurlpair=enode://{{ .Node.StrParameters.enode }}@{{ .Node.StrParameters.proxy_internal }}:30503;enode://{{ .Node.StrParameters.enode }}@{{ .Node.StrParameters.proxy_external }}:30303
--unlock={{ .Node.StrParameters.signer }}
--password=/root/.celo/configs/password.secret
--keystore=/root/.celo/configs/keystore
--ethstats={{ .Node.ID }}@baklava-ethstats.celo-testnet.org
`

	// FullnodeCmdTpl the celo command for running fullnode
	FullnodeCmdTpl = `--verbosity=2
--networkid={{ .Node.StrParameters.networkid }}
--syncmode=full
--rpc
--rpcaddr={{ .Node.StrParameters.rpcaddr }}
--rpcapi=eth,net,web3,debug,admin,personal
--light.serve={{ .Node.StrParameters.light_serve }}
--light.maxpeers={{ .Node.StrParameters.maxpeers }}
--maxpeers=1100
--etherbase={{ .Node.StrParameters.account }}
--bootnodes={{ .Node.StrParameters.bootnodes }}
--nousb
`
)