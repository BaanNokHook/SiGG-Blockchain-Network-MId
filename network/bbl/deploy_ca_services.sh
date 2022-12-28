#!/bin/bash -l

echo "Deploying CA services ..."
HOSTNAME=$HOSTNAME docker-compose -f ./peer-ca/docker-compose.yaml up --force-recreate -d postgres_ca_bbl

echo "Waiting 20 sec for database service"
sleep 20

bash ./peer-ca/set_permission.sh
bash ./peer-ca/update_config.sh

echo "Waiting 10 sec for database service"
sleep 10

HOSTNAME=$HOSTNAME docker-compose -f ./peer-ca/docker-compose.yaml up --force-recreate -d ca_bbl

echo "Finished deploying CA services ..."
