#!/bin/bash -l

docker stop $(docker ps -aq) 

docker rm $(docker ps -aq)

docker rmi $(docker images dev* -aq)

docker volume prune -f