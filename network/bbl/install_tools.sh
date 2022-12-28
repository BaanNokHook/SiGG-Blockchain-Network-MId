#!/bin/bash -l

echo "Installing tools ..."
apt-get update

apt-get install -y     apt-transport-https     ca-certificates     curl     gnupg-agent     software-properties-common     tree     zip unzip     postgresql-client
    
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -

add-apt-repository    "deb [arch=amd64] https://download.docker.com/linux/ubuntu    bionic    stable"
   
apt-get update

apt-get install docker-ce docker-ce-cli containerd.io docker-compose

snap install yq # recommend version 3.3.1
snap install jq

