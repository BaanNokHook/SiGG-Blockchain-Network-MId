#!/bin/bash -l

echo "Deploying Orderer services ..."
HOSTNAME=$HOSTNAME docker-compose -f ./orderer/docker-compose.yaml up -d
echo "Finished deploying Orderer services ..."
