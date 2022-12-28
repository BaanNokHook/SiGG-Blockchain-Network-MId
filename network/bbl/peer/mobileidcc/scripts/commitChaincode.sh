#!/bin/bash -l

peer lifecycle chaincode commit \
    -o ${ORDERER_ADDRESS} --ordererTLSHostnameOverride ${OVERIDE_ORDERER_HOSTNAME} \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
    --channelID $CHANNEL_NAME --name ${CC_NAME} \
    --peerAddresses ${CORE_PEER_ADDRESS} \
    --tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE} \
    --version ${CC_VERSION} --sequence ${CC_SEQUENCE} \
    --init-required
