#!/bin/bash -l

echo "Checking commit readiness with signature ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/checkCommitReadinessWithSignature.sh
echo "Finished commit readiness with signature checking ..."
