#!/bin/bash -l

rm -f ~/crypto-config/peerOrganizations/nbtc.mobileid.com/ca/{ca-cert.pem,tls-cert.pem}
mkdir -p ~/crypto-config/peerOrganizations/nbtc.mobileid.com/ca
while [[ ! -f ./crypto-config/peerOrganizations/nbtc.mobileid.com/ca/ca-cert.pem || ! -f ./crypto-config/peerOrganizations/nbtc.mobileid.com/ca/tls-cert.pem  ]]
do
    docker cp $(docker ps -q -f name=ca.nbtc.mobileid.com):/etc/hyperledger/fabric-ca-server/ca-cert.pem ~/crypto-config/peerOrganizations/nbtc.mobileid.com/ca/ca-cert.pem
    docker cp $(docker ps -q -f name=ca.nbtc.mobileid.com):/etc/hyperledger/fabric-ca-server/tls-cert.pem ~/crypto-config/peerOrganizations/nbtc.mobileid.com/ca/tls-cert.pem
    sleep 1
done
zip -r cainfo-nbtc.mobileid.com.zip ./crypto-config/peerOrganizations/nbtc.mobileid.com/ca/*.pem
ls -l ./crypto-config/peerOrganizations/nbtc.mobileid.com/ca/*.pem
