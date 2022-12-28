#!/bin/bash -l

# registrar operation
# CSR information and TLS MUST be configured !
# registrar should be changed for production environment.
export CA_CONNECTION=${CA_CONNECTION:=ca.bbl.mobileid.com:7054}
export CFG_PATH=~/crypto-config/peerOrganizations/bbl.mobileid.com/ca
export FABRIC_CA_CLIENT_TLS_CERTFILES=~/crypto-config/peerOrganizations/bbl.mobileid.com/ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$CFG_PATH/admin/msp

## Set default values
username=
password=
id_type=
ou_name=
affiliation=
attributes=
operation=

## Accept input values
## Override default values
while getopts ":u:p:t:o:a:r:i:" opt; do
  case $opt in
    u) username="$OPTARG"
    ;;
    p) password="$OPTARG"
    ;;
    t) id_type="$OPTARG"
    ;;
    o) ou_name="$OPTARG"
    ;;
    a) affiliation="$OPTARG"
    ;;
    r) role="$OPTARG"
    ;;
    i) operation="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done
echo "Creating new identity"
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
if [ -z "$id_type" ]
then
  read -p "Required id.type (user,admin,peer,orderer): " id_type
  echo
fi
if [ -z "$ou_name" ]
then
  read -p "Required organization unit (NBTC,AIS,BBL): " ou_name
  ou_name=$(echo $ou_name | tr '[:lower:]' '[:upper:]')
  echo
fi
if [ -z "$affiliation" ]
then
  read -p "Required affiliation (bbl.department1): " affiliation
  echo
fi
if [ "$role" == "admin" ]
then
  attributes='--id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true,abac.init=true"'
fi
if [ -z "$operation" ]
then
  read -p "Required operation: " operation
  echo
fi

# echo "${username} ${password} ${id_type} ${ou_name} ${attributes}"
# ###########################################
# # Register to RCA
# ###########################################
fabric-ca-client $operation -d \
  -u https://$CA_CONNECTION \
  --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=$ou_name

#
