#!/bin/bash -l

# registar operation
# set client's home directory which contains fabric-ca-client-config.yaml
# CSR information and TLS MUST be configured !
# registrar should be changed for production environment.
export CA_URL=${CA_URL:=localhost:7054}
export FABRIC_CA_CLIENT_TLS_CERTFILES=ca/tls-cert.pem
export FABRIC_CA_CLIENT_HOME=./crypto-config/peerOrganizations/nbtc.mobileid.com/
export OU_NAME=$(echo "orderer" | tr '[:lower:]' '[:upper:]')
mkdir -p ${FABRIC_CA_CLIENT_HOME}

fabric-ca-client enroll -d -u https://admin:adminpw@${CA_URL} --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME}

fabric-ca-client register -d --id.name peer0.nbtc.mobileid.com --id.secret peer0pw --id.type peer --id.affiliation nbtc.department1  -u https://admin:adminpw@${CA_URL}  --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME}
fabric-ca-client register -d --id.name peer1.nbtc.mobileid.com --id.secret peer1pw --id.type peer --id.affiliation nbtc.department1  -u https://admin:adminpw@${CA_URL}  --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME}

fabric-ca-client register -d --id.name Admin@nbtc.mobileid.com --id.secret adminpw --id.type admin --id.affiliation nbtc.department1 --id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true,abac.init=true" -u https://admin:adminpw@${CA_URL}  --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME}
fabric-ca-client register -d --id.name User1@nbtc.mobileid.com --id.secret user1pw --id.type user --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME}

fabric-ca-client identity list
tree ./crypto-config/
