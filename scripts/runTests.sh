#!/bin/bash

NODES=(${@:2})
BINARY=$1
PROJECT_ROOT="$( cd "$( dirname "$0" )/.." && pwd )"

echo ""
echo "Blockdaemon BPM Celo Integration Test Suite"
echo ""
echo "Project Root: ${PROJECT_ROOT}"
echo "binary: ${1}"
echo "testing: ${NODES}"
echo ""

function main() {

    for subtype in "${NODES[@]}"; do

        echo ""
        echo "Creating ${subtype}"
        cd $PROJECT_ROOT/build/$subtype
        clean $subtype

        setupWorkspace $subtype
        configureSubtype $subtype
        startSubtype $subtype
        echo "finished ${subtype}!"
        echo ""
    done

    for subtype in "${NODES[@]}"; do

        if [ "$subtype" != "attestation-service" ]; then
            echo ""
            echo "Testing $subtype"
            cd $PROJECT_ROOT/build/$subtype
            ./$BINARY test node.$subtype.json
            echo "finished ${subtype}!"
            echo ""
        fi
    done
}

function setupWorkspace() {
    echo "scaffolding: [$1]"
    mkdir $PROJECT_ROOT/build/$1/logs $PROJECT_ROOT/build/$1/monitoring 2> /dev/null
    touch $PROJECT_ROOT/build/$1/monitoring/filebeat.yml 2> /dev/null
    cp $PROJECT_ROOT/build/$BINARY .

    # if validator, setup validator
    if [ "$1" == "validator" ]; then
        setupValidator $1
    fi
}

function setupValidator() {

    echo "updating validator json with proxy details..."
    proxy=bpm-celo-test-proxy

    # get enode url, ip
    enode=$(docker exec ${proxy} geth --exec "admin.nodeInfo['enode'].split('//')[1].split('@')[0]" attach | tr -d '"')
    proxyIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${proxy})
    externalIP=$(dig +short myip.opendns.com @resolver1.opendns.com)
    echo "proxy enode: ${enode}"
    echo "proxy ip: ${proxyIP}"

    cat <<< $(jq --arg enode "$enode" '.str_parameters.enode = $enode' node.validator.json) > node.validator.json
    cat <<< $(jq --arg externalIP "$externalIP" '.str_parameters.proxy_external = $externalIP' node.validator.json) > node.validator.json
    cat <<< $(jq --arg proxyIP "$proxyIP" '.str_parameters.proxy_internal = $proxyIP' node.validator.json) > node.validator.json

    echo "... done updating validator json"
}

function configureSubtype() {
    echo "configuring: [$1]"
    echo $PWD
    ./$BINARY create-configurations node.$1.json
}

function startSubtype() {
    echo "starting: [$1]"
    ./$BINARY start node.$1.json
}

# TODO: implement this somewhere, currently nodes will need to stay up for rpc testing.
function destroyAll() {
    docker rm $(docker ps --filter name=bpm-celo-test*)
    exit
}

function clean() {
    if [ "$1" != "" ]; then
        echo "cleaning up [$1]..."
        ls | grep -v node.$1.json | xargs rm -rf
    fi
}

main
