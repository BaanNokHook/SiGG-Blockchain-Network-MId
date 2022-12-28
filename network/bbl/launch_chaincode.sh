#!/bin/bash -l

# for Fabric version 1.4
echo "Rebuild chaincode instance ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) go build
echo "Run dry query to launch chaincode instance ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/queryChaincode.sh
echo "List running containers, please verify that the new chaincode instance is up and running ..."
docker images
docker ps
