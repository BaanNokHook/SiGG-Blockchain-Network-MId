#!/bin/bash -l
peer lifecycle chaincode queryinstalled
peer lifecycle chaincode queryinstalled > log.txt
export PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
echo $PACKAGE_ID
