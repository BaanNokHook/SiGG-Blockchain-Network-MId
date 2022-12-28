echo "Deploying Peer services ..."
HOSTNAME=$HOSTNAME docker-compose -f ./peer/docker-compose.yaml up -d

echo "Waiting for Cli starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=cli_peer))" == false ]; do sleep 1; done

echo "Finished deploying Peer services ..."
