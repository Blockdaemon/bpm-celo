# Celo BPM Plugin

This plugin is for the [Blockdaemon Protocol Manager](https://gitlab.com/Blockdaemon/bpm-cli).

Installation
```
bpm packages install celo
```

## Validator
(depends on a proxy node being available first)

```
bpm nodes configure celo
    --subtype=validator
    --networkid
    --proxy_external
    --proxy_internal
    --signer
```

## Proxy

```
bpm nodes configure celo
    --subtype=proxy
    --networkid
    --signer
```

## Fullnode

```
bpm nodes configure celo
    --subtype=fullnode
    --networkid
    --signer
```

## Accounts

```
bpm nodes configure celo
    --subtype=accounts
    --networkid
```

## Attestations
(creates 2 nodes and 1 postgres container)

```
bpm nodes configure celo
    --subtype=attestations
    --networkid
    --signer
    --validator
    --twillio-account-sid
    --twillio-service-sid
    --twillio-auth-token
```
