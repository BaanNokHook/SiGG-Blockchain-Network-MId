#!/bin/bash -l

echo "Checking chaincode committed ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/queryCommitted.sh
echo "Finished chaincode committed checking ..."
