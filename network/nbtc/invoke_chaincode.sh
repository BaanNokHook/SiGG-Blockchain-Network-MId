#!/bin/bash -l

echo "Waiting for cli starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com))" == false ]; do sleep 1; done

echo "Invoking chaincode ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/invokeChaincode.sh
echo "Finished chaincode invoking ..."

