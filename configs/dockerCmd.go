package configs

const (
	// CollectorEnvTpl tempalte for collector command
	CollectorEnvTpl = `SERVICE_PORT=8545
SERVICE_HOST=bpm-{{ .Node.ID }}-{{ .Node.StrParameters.subtype }}`

	// ProxyCmdTpl the celo command for running proxies
	ProxyCmdTpl = `--verbosity=3
--networkid={{ .Node.StrParameters.networkid }}
--syncmode=full
--proxy.proxy
--proxy.proxiedvalidatoraddress={{ .Node.StrParameters.signer }}
--proxy.internalendpoint=:30503
--rpcvhosts=bpm-{{ .Node.ID }}-{{ .Node.StrParameters.subtype }}
--etherbase={{ .Node.StrParameters.signer }}
--bootnodes={{ .Node.StrParameters.bootnodes }}
<<<<<<< HEAD
{{ if index .Node.StrParameters.celo }}{{ .Node.StrParameters.celo }}{{ end }}
=======
>>>>>>> release-0.0.2
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
--rpcvhosts=bpm-{{ .Node.ID }}-{{ .Node.StrParameters.subtype }}
<<<<<<< HEAD
--password=/root/.celo/configs/password.secret
--keystore=/root/.celo/configs/keystore
{{ if index .Node.StrParameters.celo }}{{ .Node.StrParameters.celo }}{{ end }}
=======
--password=/root/.celo/configs/.password.secret
--keystore=/root/.celo/configs/keystore
>>>>>>> release-0.0.2
`

	// FullnodeCmdTpl the celo command for running fullnode
	FullnodeCmdTpl = `--verbosity=3
--networkid={{ .Node.StrParameters.networkid }}
--syncmode=full
--rpc
--rpcaddr={{ .Node.StrParameters.rpcaddr }}
--rpcapi=eth,net,web3,debug,admin,personal
--light.serve={{ .Node.StrParameters.light_serve }}
--light.maxpeers={{ .Node.StrParameters.light_maxpeers }}
--maxpeers={{ .Node.StrParameters.maxpeers }}
--rpcvhosts=bpm-{{ .Node.ID }}-{{ .Node.StrParameters.subtype }}
--etherbase={{ .Node.StrParameters.account }}
--bootnodes={{ .Node.StrParameters.bootnodes }}
<<<<<<< HEAD
{{ if index .Node.StrParameters.celo }}{{ .Node.StrParameters.celo }}{{ end }}
`
=======
{{ if eq .Node.StrParameters.nousb "true" "TRUE" "True" }}--nousb{{ end }}
`

	// AttestationCmdTpl the celo command for running attestation node
	AttestationCmdTpl = `--verbosity=3
--networkid={{ .Node.StrParameters.networkid }}
--syncmode=full
--rpc
--rpcvhosts=bpm-{{ .Node.ID }}-{{ .Node.StrParameters.subtype }}
--rpcaddr={{ .Node.StrParameters.rpcaddr }}
--rpcapi=eth,net,web3,debug,admin,personal
--allow-insecure-unlock
--unlock={{ .Node.StrParameters.signer }}
--keystore=/root/.celo/configs/keystore
--password=/root/.celo/configs/.password.secret
--bootnodes={{ .Node.StrParameters.bootnodes }}
`

	// AttestationServiceCmdTpl the celo command for running attestation service
	AttestationServiceCmdTpl = ``
>>>>>>> release-0.0.2
)
