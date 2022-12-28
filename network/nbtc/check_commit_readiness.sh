#!/bin/bash -l

echo "Checking commit readiness ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/checkCommitReadiness.sh
echo "Finished commit readiness checking ..."
