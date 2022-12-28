#!/bin/bash -l
peer lifecycle chaincode queryinstalled
peer lifecycle chaincode queryinstalled > log.txt
export PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
echo "Approving Chaincode with Package ID ${PACKAGE_ID} seq. ${CC_SEQUENCE}"

peer lifecycle chaincode approveformyorg \
    -o ${ORDERER_ADDRESS} --ordererTLSHostnameOverride ${OVERIDE_ORDERER_HOSTNAME} \
    --name "${CC_NAME}" \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
    --channelID $CHANNEL_NAME \
    --version ${CC_VERSION} --init-required \
    --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} 
