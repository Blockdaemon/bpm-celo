# Celo BPM Plugin

Celo             |  Blockdaemon
:-------------------------:|:-------------------------:
![](/assets/celo-logo-200x200.png)  |  ![](/assets/bd-logo-200x200.png)

This plugin is for the [Blockdaemon Protocol Manager](https://gitlab.com/Blockdaemon/bpm-cli).

Installation
```
bpm packages install celo
```

## Arguments

### Validator
(depends on a proxy node being available first)

```
bpm nodes configure celo
    --subtype=validator
    --networkid
    --proxy_external
    --proxy_internal
    --proxy_enode
    --port
    --signer
    --keystore-file
    --keystore-pass
```

### Proxy

```
bpm nodes configure celo
    --subtype=proxy
    --networkid
    --signer
    --bootnodes
```

### Fullnode

```
bpm nodes configure celo
    --subtype=fullnode
    --networkid
    --account
    --bootnodes
```

## Running Nodes

### Validator and Proxy

In order to run a validator you must run a proxy first. To run a proxy call:
```
bpm nodes configure celo
    --subtype=proxy
    --networkid=40120
    --signer=0xf2334aae1b2f273b600abff9a491eb720d842b6d
    --bootnodes=enode://5aaf10664b12431c250597e980aacd7d5373cae00f128be5b00364344bb96bce7555b50973664bddebd1cb7a6d3fb927bec81527f80e22a26fa373c375fcdefc@35.247.75.229:30301
```

You can get the `bootnodes` parameter running the following:
(replace `$CELO_IMAGE` with the image for the relevant change, eg `us.gcr.io/celo-testnet/celo-node:baklava`)
```
docker run --rm --entrypoint cat $CELO_IMAGE /celo/bootnodes"
```

Start the proxy
```
bpm nodes start bpm-celo-proxy
```

Once running you need the `enode url` and `public ip` so that your validator can
connect to it. Run the following:
```
docker exec celo-proxy geth --exec "admin.nodeInfo['enode'].split('//')[1].split('@')[0]" attach | tr -d '"'
dig +short myip.opendns.com @resolver1.opendns.com
```

Next we can configure our validator:
```
bpm nodes configure celo
    --subtype=validator
    --networkid=40120
    --signer=0xf2334aae1b2f273b600abff9a491eb720d842b6d
    --port=30303
    --proxy_internal=<proxy internal ip, if none then use external ip from previous command>
    --proxy_external=<proxy external ip from previous command>
    --enode=<proxy enode string from previous command>
    --keystore-file=/path/to/signer/keystoreJSON
    --keystore-pass=/path/to/keystore/password.secret


```

If running on same instance as proxy make sure you change the listening port on
the validator to something other than `30303` as proxy needs to to communicate
with the interweb.

### Fullnode

A fullnode can be run using the following:
```
bpm --debug nodes configure celo --subtype=fullnode --networkid=40120 --account=0xf2334aae1b2f273b600abff9a491eb720d842b6d --port=30314 --bootnodes=enode://5aaf10664b12431c250597e980aacd7d5373cae00f128be5b00364344bb96bce7555b50973664bddebd1cb7a6d3fb927bec81527f80e22a26fa373c375fcdefc@34.82.45.71:30301
```

## Development

To develop with this plugin.

First generate the config files:
```
go run ./cmd/main.go create-configurations node.json
```

This will create the following files in the repository:
```
...
├── celo.dockercmd
├── configs
│   ├── collector.env
├── node.json
...
```

When the plugin is called by `bpm-cli` there is the default location for
configurations and data, usually `~/.bpm`. But here there will be none as we
are calling from within the plugin. So we need to create a `node.json` file:

`node.json`
```
{
  "id": "celo",
  "plugin": "celo",
  "str_parameters": {
    "docker-network": "bpm",
    "subtype": "validator",
    "networkid": "40120",
    "signer": "0xc1c048DE906CE7e3F99c1feC0651671ec91970F9",
    "proxy_internal": "0.0.0.0",
    "proxy_external": "1.1.1.1",
    "proxy_enode": "5aaf10664b12431c250597e980aacd7d5373cae00f128be5b00364344bb96bce7555b50973664bddebd1cb7a6d3fb927bec81527f80e22a26fa373c375fcdefc",
    "bootnodes": "enode://5aaf10664b12431c250597e980aacd7d5373cae00f128be5b00364344bb96bce7555b50973664bddebd1cb7a6d3fb927bec81527f80e22a26fa373c375fcdefc@34.82.45.71:30301"
  },
  "version": "1.0.0"
}
```

Then run:
```
go run ./cmd/main.go start node.json
```
