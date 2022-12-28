#!/bin/bash -l

HOSTNAME=$HOSTNAME docker-compose -f ./orderer/docker-compose.yaml down
HOSTNAME=$HOSTNAME docker-compose -f ./peer/docker-compose.yaml down
HOSTNAME=$HOSTNAME docker-compose -f ./peer-ca/docker-compose.yaml down
HOSTNAME=$HOSTNAME docker-compose -f ./middleware/docker-compose.yaml down
