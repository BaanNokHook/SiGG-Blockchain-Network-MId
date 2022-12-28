#!/bin/bash -l

echo "Setting host ..."
USERNAME=$1
HOSTNAME=$2

hostnamectl set-hostname ${HOSTNAME}

groupadd -f docker
usermod -aG docker $USERNAME
