#!/bin/bash -l

cd && openssl rand -out .rnd &

cd chaincode && ./build.sh &

cd network/ais/middleware && ./build.sh &

docker network create -d bridge mobileid-network &

wait