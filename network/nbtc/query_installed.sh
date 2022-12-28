#!/bin/bash -l

echo "Checking chaincode installed ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/queryInstalled.sh
echo "Finished chaincode installed checking ..."
