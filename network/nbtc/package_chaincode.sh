#!/bin/bash -l

echo "Package chaincode ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/packageChaincode.sh
echo "Finished package chaincode ..."
