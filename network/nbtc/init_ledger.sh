#!/bin/bash -l

echo "init ledger ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.nbtc.mobileid.com) bash scripts/initLedger.sh
echo "Finished init ledger ..."
