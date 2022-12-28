#!/bin/bash -l

# registrar operation
# CSR information and TLS MUST be configured !
# registrar should be changed for production environment.
export CA_CONNECTION=${CA_CONNECTION:=ca.bbl.mobileid.com:7054}
export CFG_PATH=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$CFG_PATH/admin/msp

## Set default values
username=
password=

## Accept input values
## Override default values
while getopts ":u:p:t:o:a:r:" opt; do
  case $opt in
    u) username="$OPTARG"
    ;;
    p) password="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done
echo "Enroll admin identity"
if [ -z "$username" ]
then
  read -p "Required username (cn): " username
  echo
fi
if [ -z "$password" ]
then
  read -s -p "Required password: " password
  echo
fi

# ###########################################
# Enroll Admin identity for CA administration
# ###########################################
fabric-ca-client enroll -d \
  -u https://${username}:${password}@${CA_CONNECTION} \
  --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=$ou_name


#
