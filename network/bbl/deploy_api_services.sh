#!/bin/bash -l

echo "Deploying API services ..."
HOSTNAME=$HOSTNAME docker-compose -f ./middleware/docker-compose.yaml up -d --build --force-recreate
# if fails related to API server, please clean and use:
# HOSTNAME=$HOSTNAME docker-compose up -d --build --no-cache
echo "Finished deploying API services ..."
