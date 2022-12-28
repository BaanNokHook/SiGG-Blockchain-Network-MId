#!/bin/bash -l

echo "Waiting for database server starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=postgres.nbtc.mobileid.com))" == false ]; do sleep 1; done

docker exec -u postgres -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) cp /config/postgresql.conf /data/postgres/
docker exec -u postgres -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) cp /config/pg_hba.conf /data/postgres/
docker exec -u postgres -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) cp /config/pg_ident.conf /data/postgres/
docker exec -u postgres -it $(docker ps -q -f name=postgres.nbtc.mobileid.com) pg_ctl reload

echo "Updated config in database server with gracefully reload"
