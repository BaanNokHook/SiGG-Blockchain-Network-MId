#!/bin/bash -l

echo "Installing chaincode ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/installChaincode.sh
echo "Chaincode installed ..."
