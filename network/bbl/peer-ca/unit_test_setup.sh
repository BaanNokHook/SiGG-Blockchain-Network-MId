#!/bin/bash -l

echo "Waiting for ca server starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=ca.*.mobileid.com))" == false ]; do sleep 1; done

echo " echo \"Install unit testing tools\"
  apt-get update
  apt-get upgrade -y
  apt-get install -y postgresql-client
" | docker exec -i $(docker ps -q -f name=ca.*.mobileid.com) bash --
