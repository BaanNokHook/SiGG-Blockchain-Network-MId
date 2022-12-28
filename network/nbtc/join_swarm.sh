#!/bin/bash -l

echo "Joining Docker Swarm ..."
HOST_IP=$1
SWARM_TOKEN=$2
SWARM_LEADER=$3

function joinSwarmIfNeed(){
if [ "$(docker info | grep Swarm | sed 's/Swarm: //g')" == "inactive" ]; then
    docker swarm join --token ${SWARM_TOKEN} ${SWARM_LEADER}:2377 --advertise-addr ${HOST_IP}
else
    echo "This instance is in Swarm (please check it is the same swarm as another instances)"
fi
}
joinSwarmIfNeed

