#!/bin/bash -l

echo "Waiting for database server starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=postgres.nbtc.mobileid.com))" == false ]; do sleep 1; done

docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chown postgres /var/lib/postgresql/server.pem
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chown postgres /var/lib/postgresql/server.key
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chown postgres /var/lib/postgresql/root.pem
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chown postgres /var/lib/postgresql/root.key
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chmod 0600 /var/lib/postgresql/server.pem
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chmod 0600 /var/lib/postgresql/server.key
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chmod 0600 /var/lib/postgresql/root.pem
docker exec -u 0 -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) chmod 0600 /var/lib/postgresql/root.key
echo "Changed permission of server cert files"
