#!/bin/bash -l

rm -f ~/crypto-config/peerOrganizations/bbl.mobileid.com/ca/{ca-cert.pem,tls-cert.pem}
mkdir -p ~/crypto-config/peerOrganizations/bbl.mobileid.com/ca
while [[ ! -f ./crypto-config/peerOrganizations/bbl.mobileid.com/ca/ca-cert.pem || ! -f ./crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem  ]]
do
    docker cp $(docker ps -q -f name=ca.bbl.mobileid.com):/etc/hyperledger/fabric-ca-server/ca-cert.pem ~/crypto-config/peerOrganizations/bbl.mobileid.com/ca/ca-cert.pem
    docker cp $(docker ps -q -f name=ca.bbl.mobileid.com):/etc/hyperledger/fabric-ca-server/tls-cert.pem ~/crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem
    sleep 1
done
zip -r cainfo-bbl.mobileid.com.zip ./crypto-config/peerOrganizations/bbl.mobileid.com/ca/*.pem
ls -l ./crypto-config/peerOrganizations/bbl.mobileid.com/ca/*.pem
