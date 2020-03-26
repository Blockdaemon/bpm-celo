package configs

const (
	// CollectorEnvTpl tempalte for collector command
	CollectorEnvTpl = `SERVICE_PORT=9933
SERVICE_HOST={{ .Node.NamePrefix }}proxy`

	// CeloCmdTpl the celo command
	CeloCmdTpl = `
--verbosity 3
--networkid {{ .Node.StrParameters.networkid }}
--syncmode full
--proxy.proxy
--proxy.proxiedvalidatoraddress {{ .Node.StrParameters.signer }}
--proxy.internalendpoint :30503
--etherbase {{ .Node.StrParameters.signer }}
--bootnodes {{ .Node.StrParameters.bootnodes }}
--ethstats={{ .Node.ID }}-proxy@baklava-ethstats.celo-testnet.org
`
)
