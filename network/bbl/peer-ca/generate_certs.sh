#!/bin/bash -l

###########################################
# Setup RCA
###########################################

export URL_TLSCA=${URL_TLSCA:=ca.mobileid.com:7054}
export URL_ORG_CA=${URL_ORG_CA:=ca.bbl.mobileid.com:7054}

export FABRIC_CA_CLIENT_MSPDIR=
export FABRIC_CA_CLIENT_TLS_CERTFILES=
export FABRIC_CA_CLIENT_HOME=
export CFG_PATH=$PWD/crypto-config
export OU_NAME=$(echo "bbl" | tr '[:lower:]' '[:upper:]')

export CRYPTO_CONFIG_PATH=$PWD/crypto-config

mkdir -p $CFG_PATH/
mkdir -p crypto-config/peerOrganizations/bbl.mobileid.com/ca/admin/msp/{cacerts,keystore,signcerts,user}
mkdir -p crypto-config/tlsca/admin/msp/{cacerts,keystore,signcerts,user}

###########################################
# Setup Peer0 For bbl
###########################################

export FABRIC_CA_CLIENT_MSPDIR=
export FABRIC_CA_CLIENT_TLS_CERTFILES=
export FABRIC_CA_CLIENT_HOME=

mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/{ca,tlsca}
mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}
mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/
mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}
mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}
mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/msp/{admincerts,cacerts,tlscacerts}

# sudo chmod -R 777 ./

# Enrollment cert crypto
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem
export FABRIC_CA_CLIENT_HOME=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com
export FABRIC_CA_CLIENT_MSPDIR=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp
fabric-ca-client enroll --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME} -u https://peer0.bbl.mobileid.com:peer0pw@$URL_ORG_CA

# TLS crypto
export FABRIC_CA_CLIENT_MSPDIR=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-tmp/tls-cert.pem
fabric-ca-client enroll --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME} -u https://peer0.bbl.mobileid.com:peer0pw@$URL_TLSCA --enrollment.profile tls --csr.hosts peer0.bbl.mobileid.com
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/keystore/*_sk | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/keystore/key.pem

# Admin
export FABRIC_CA_CLIENT_MSPDIR=
export FABRIC_CA_CLIENT_TLS_CERTFILES=
export FABRIC_CA_CLIENT_HOME=

export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem
export FABRIC_CA_CLIENT_HOME=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com
export FABRIC_CA_CLIENT_MSPDIR=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp
fabric-ca-client enroll --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME} -u https://Admin@bbl.mobileid.com:adminpw@$URL_ORG_CA

cp $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/signcerts/cert.pem $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/admincerts/bbl-admin-cert.pem

# user1
export FABRIC_CA_CLIENT_MSPDIR=
export FABRIC_CA_CLIENT_TLS_CERTFILES=
export FABRIC_CA_CLIENT_HOME=

export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem
export FABRIC_CA_CLIENT_HOME=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com
export FABRIC_CA_CLIENT_MSPDIR=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp
fabric-ca-client enroll --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME} -u https://User1@bbl.mobileid.com:user1pw@$URL_ORG_CA

###########################################
# Setup Peer1 For bbl
###########################################

export FABRIC_CA_CLIENT_MSPDIR=
export FABRIC_CA_CLIENT_TLS_CERTFILES=
export FABRIC_CA_CLIENT_HOME=

mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/{admincerts,cacerts,keystore,signcerts,tlscacerts}
mkdir -p $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/

# sudo chmod -R 777 ./

# Enrollment cert crypto
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem
export FABRIC_CA_CLIENT_HOME=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com
export FABRIC_CA_CLIENT_MSPDIR=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp
fabric-ca-client enroll --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME} -u https://peer1.bbl.mobileid.com:peer1pw@$URL_ORG_CA

# TLS crypto
export FABRIC_CA_CLIENT_MSPDIR=$CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PWD/crypto-tmp/tls-cert.pem
fabric-ca-client enroll --csr.names C=TH,ST=Bangkok,L="Phaya Thai",O=MobileID,OU=${OU_NAME} -u https://peer1.bbl.mobileid.com:peer1pw@$URL_TLSCA --enrollment.profile tls --csr.hosts peer1.bbl.mobileid.com
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/keystore/*_sk | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/keystore/key.pem

cp $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/signcerts/cert.pem $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/admincerts/bbl-admin-cert.pem

###########################################
# Move Cert
###########################################

export FABRIC_CA_CLIENT_MSPDIR=
export FABRIC_CA_CLIENT_TLS_CERTFILES=
export FABRIC_CA_CLIENT_HOME=
export CFG_PATH=$PWD/crypto-config

# folder users - Admin - msp
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/tlscacerts/tlsca.bbl.mobileid.com-cert.pem
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/admincerts/Admin@bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/cacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/cacerts/ca.bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/keystore/* | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/keystore/Admin@bbl.mobileid.com-key.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/signcerts/Admin@bbl.mobileid.com-cert.pem

# folder users - user1 - msp
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/tlscacerts/tlsca.bbl.mobileid.com-cert.pem
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/admincerts/User1@bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/cacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/cacerts/ca.bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/keystore/* | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/keystore/User1@bbl.mobileid.com-key.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/User1@bbl.mobileid.com/msp/signcerts/User1@bbl.mobileid.com-cert.pem

# folder msp
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/users/Admin@bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/msp/admincerts/Admin@bbl.mobileid.com-cert.pem
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/cacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/msp/cacerts/ca.bbl.mobileid.com-cert.pem
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/msp/tlscacerts/tlsca.bbl.mobileid.com-cert.pem

# folder tlsca
CERT_FILE_PATH=$PWD/crypto-config/peerOrganizations/bbl.mobileid.com/ca/tls-cert.pem
cp $CERT_FILE_PATH $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/tlsca/tlsca.bbl.mobileid.com-cert.pem

# peer0
# folder peers - tls
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/ca.crt
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/server.crt
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/keystore/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/server.key

# folder peers - msp
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/tlscacerts/tlsca.bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/admincerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/admincerts/Admin@bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/cacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/cacerts/ca.bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/keystore/* | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/keystore/peer0.bbl.mobileid.com-key.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/msp/signcerts/peer0.bbl.mobileid.com-cert.pem

# peer1
# folder peers - tls
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/ca.crt
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/server.crt
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/keystore/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/server.key

# folder peers - msp
cp $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/tlscacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/tlscacerts/tlsca.bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/admincerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/admincerts/Admin@bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/cacerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/cacerts/ca.bbl.mobileid.com-cert.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/keystore/* | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/keystore/peer1.bbl.mobileid.com-key.pem
mv $(ls $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/signcerts/*.pem | head -n1) $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/msp/signcerts/peer1.bbl.mobileid.com-cert.pem

rm -rf $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/{signcerts,tlscacerts,keystore,user,cacerts}
rm -rf $CRYPTO_CONFIG_PATH/peerOrganizations/bbl.mobileid.com/peers/peer1.bbl.mobileid.com/tls/{signcerts,tlscacerts,keystore,user,cacerts}

find crypto-config -maxdepth 10 -type f
(
zip -r ~/bbl.mobileid.com-msp.zip ./crypto-config/peerOrganizations/bbl.mobileid.com/msp/*/*.pem
)
