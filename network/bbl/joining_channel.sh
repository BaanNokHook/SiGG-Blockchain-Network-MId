#!/bin/bash -l

echo "Joining channel ..."
echo "Waiting for cli starting up ..."
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/joinChannel.sh
docker exec -t $(docker ps -q -f name=cli_peer0.bbl.mobileid.com) bash scripts/updateAnchorPeer.sh
echo "Finished channel joining ..."
