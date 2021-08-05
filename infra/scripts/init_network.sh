#!/bin/bash

CWD=$(pwd)

trap "finish" INT TERM

NETWORK_INSPECTION=$(docker network inspect "atai_envoy")
EXITCODE_NETWORK_INSPECTION=$?
if [ $EXITCODE_NETWORK_INSPECTION -ne 0 ]
then
    echo "create new network"
    docker network create \
        --driver="bridge" \
        --subnet="172.18.0.0/24" \
        --gateway="172.18.0.1" \
        atai_envoy
else
    echo "network already exists"
fi

finish() {
    local existcode=$?
    cd $CWD
    exit $existcode
}