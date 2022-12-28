#!/bin/bash -l

echo "Commit chaincode ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/commitChaincode.sh
echo "Finished chaincode commited ..."
