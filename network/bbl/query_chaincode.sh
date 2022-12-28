#!/bin/bash -l

echo "Querying chaincode ..."
echo "Waiting for cli starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=cli_peer0.bbl.mobileid.com))" == false ]; do sleep 1; done

docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/queryChaincode.sh
echo "Chaincode queried ..."
