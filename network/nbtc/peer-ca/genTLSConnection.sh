#!/bin/bash -l

# https://www.makethenmakeinstall.com/2014/05/ssl-client-authentication-step-by-step/

# keyrequest:
#   algo: rsa
#   size: 4096
# names:
#    - C: TH
#      ST: "Bangkok"
#      L:
#      O: MobileID
#      OU: NBTC


OU_NAME=$(echo ${OU_NAME:=nbtc} | tr '[:lower:]' '[:upper:]')
C=TH
ST=Bangkok
O=MobileID
CN=nbtc.mobileid.com

SERVER_DIR=${SERVER_DIR:=${HOME}/peer-ca/postgres}
CLIENT_DIR=${CLIENT_DIR:=${HOME}/peer-ca/fabric-ca}
TEST_DIR=${TEST_DIR:=${HOME}/peer-ca/test}

# use for test
SERVER_DIR=./peer-ca/postgres
CLIENT_DIR=./peer-ca/fabric-ca
TEST_DIR=./peer-ca/test

mkdir -p ${SERVER_DIR}
mkdir -p ${CLIENT_DIR}
mkdir -p ${TEST_DIR}/postgres
mkdir -p ${TEST_DIR}/fabric-ca

echo "####################"
echo "Generating root cert"
echo "####################"
openssl genrsa -out ${SERVER_DIR}/root.key 4096
openssl req -new -key ${SERVER_DIR}/root.key \
  -subj "/C=${C}/ST=${ST}/O=${O}/OU=${OU_NAME}/CN=${CN}" \
  -out ${SERVER_DIR}/root.req -sha256
openssl x509 -req -in ${SERVER_DIR}/root.req -signkey ${SERVER_DIR}/root.key \
  -set_serial 101 -extensions server -days 3650 -outform PEM -out ${SERVER_DIR}/root.pem -sha256
chmod 600 ${SERVER_DIR}/root.pem
chmod 600 ${SERVER_DIR}/root.key

echo "######################"
echo "Generating server cert"
echo "######################"
openssl genrsa -out ${SERVER_DIR}/server.key 4096
openssl req -new -key ${SERVER_DIR}/server.key \
  -subj "/C=${C}/ST=${ST}/O=${O}/OU=${OU_NAME}/CN=postgres.${CN}" \
  -out ${SERVER_DIR}/server.req -sha256
openssl x509 -req -in ${SERVER_DIR}/server.req -CA ${SERVER_DIR}/root.pem \
  -CAkey ${SERVER_DIR}/root.key \
  -set_serial 101 -extensions server -days 3650 -outform PEM -out ${SERVER_DIR}/server.pem -sha256
chmod 600 ${SERVER_DIR}/server.pem
chmod 600 ${SERVER_DIR}/server.key

echo "######################"
echo "Generating client cert"
echo "######################"
openssl genrsa -out ${CLIENT_DIR}/client.key 4096
openssl req -new -key ${CLIENT_DIR}/client.key \
  -subj "/C=${C}/ST=${ST}/O=${O}/OU=${OU_NAME}/CN=ca.${CN}" \
  -out ${CLIENT_DIR}/client.req -sha256
openssl x509 -req -in ${CLIENT_DIR}/client.req -CA ${SERVER_DIR}/root.pem \
  -CAkey ${SERVER_DIR}/root.key \
  -set_serial 101 -extensions client -days 3650 -outform PEM -out ${CLIENT_DIR}/client.pem
chmod 600 ${CLIENT_DIR}/client.pem
chmod 600 ${CLIENT_DIR}/client.key
cp ${SERVER_DIR}/root.pem ${CLIENT_DIR}
chmod 600 ${CLIENT_DIR}/root.pem

# rm -f $SERVER_DIR/server.key $SERVER_DIR/server.pem $SERVER_DIR/server.req
# rm -f $CLIENT_DIR/client.key $CLIENT_DIR/client.pem $CLIENT_DIR/client.req

## Generate test certs
##

openssl genrsa -out ${TEST_DIR}/postgres/intruder_root.key 4096
openssl req -new -key ${TEST_DIR}/postgres/intruder_root.key \
  -subj "/C=${C}/ST=${ST}/O=${O}/OU=${OU_NAME}/CN=${CN}" \
  -out ${TEST_DIR}/postgres/intruder_root.req -sha256
openssl x509 -req -in ${TEST_DIR}/postgres/intruder_root.req -signkey ${TEST_DIR}/postgres/intruder_root.key \
  -set_serial 101 -extensions server -days 3650 -outform PEM -out ${TEST_DIR}/postgres/intruder_root.pem -sha256
chmod 600 ${TEST_DIR}/postgres/intruder_root.pem
chmod 600 ${TEST_DIR}/postgres/intruder_root.key

openssl genrsa -out ${TEST_DIR}/fabric-ca/intruder_client.key 4096
openssl req -new -key ${TEST_DIR}/fabric-ca/intruder_client.key \
  -subj "/C=${C}/ST=${ST}/O=${O}/OU=${OU_NAME}/CN=ca.${CN}" \
  -out ${TEST_DIR}/fabric-ca/intruder_client.req -sha256
openssl x509 -req -in ${TEST_DIR}/fabric-ca/intruder_client.req -CA ${TEST_DIR}/postgres/intruder_root.pem \
  -CAkey ${TEST_DIR}/postgres/intruder_root.key \
  -set_serial 101 -extensions client -days 3650 -outform PEM -out ${TEST_DIR}/fabric-ca/intruder_client.pem
chmod 600 ${TEST_DIR}/fabric-ca/intruder_client.pem
chmod 600 ${TEST_DIR}/fabric-ca/intruder_client.key

# rm -f $SERVER_DIR/server.key $SERVER_DIR/server.pem $SERVER_DIR/server.req
# rm -f $CLIENT_DIR/client.key $CLIENT_DIR/client.pem $CLIENT_DIR/client.req

# Test generated certs with openssl
echo "##########################"
echo "Verify server: expected OK"
echo "##########################"
openssl verify -CAfile ${SERVER_DIR}/root.pem \
  ${SERVER_DIR}/server.pem

echo "##########################"
echo "Verify client: expected OK"
echo "##########################"
openssl verify -CAfile ${SERVER_DIR}/root.pem \
  ${CLIENT_DIR}/client.pem

# Test with mocking intruder
echo "########################################"
echo "Verify intruder with its ca: expected OK"
echo "########################################"
openssl verify -CAfile ${TEST_DIR}/postgres/intruder_root.pem \
  ${TEST_DIR}/fabric-ca/intruder_client.pem

echo "#########################################################"
echo "Verify intruder with cross ca: expected signature failure"
echo "#########################################################"
openssl verify -CAfile ${SERVER_DIR}/root.pem \
  ${TEST_DIR}/fabric-ca/intruder_client.pem
