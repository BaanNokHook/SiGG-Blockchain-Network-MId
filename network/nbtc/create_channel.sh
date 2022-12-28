#!/bin/bash -l

echo "Create channel ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/createChannel.sh
echo "Finished channel created ..."
