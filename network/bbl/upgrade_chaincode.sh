#!/bin/bash -l

# for Fabric version 1.4
CC_VERSION=bbl

echo "Updating cli for new configs ..."
docker-compose -f peer/docker-compose.yaml up --force-recreate -d

echo "Waiting for cli starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=cli_peer0.bbl.mobileid.com))" == false ]; do sleep 1; done

echo "Installing new chaincode ..."
docker exec -e CC_VERSION= -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/installChaincode.sh

echo "Upgrading new chaincode ..."
docker exec -e CC_VERSION= -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/upgradeChaincode.sh

echo "Please verify chaincode digest ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) peer chaincode list --installed
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) peer chaincode list --instantiated -C mobileid

# echo "Then, Channel Admin must mannually instantiate new chaincode"
