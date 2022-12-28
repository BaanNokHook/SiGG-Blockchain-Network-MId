#!/bin/bash -l

echo "Approve chaincode ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/approveChaincode.sh
echo "Finished chaincode approve ..."
