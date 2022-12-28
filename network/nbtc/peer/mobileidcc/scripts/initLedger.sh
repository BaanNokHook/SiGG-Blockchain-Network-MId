#!/bin/bash -l

peer chaincode invoke \
    -o ${ORDERER_ADDRESS} --ordererTLSHostnameOverride ${OVERIDE_ORDERER_HOSTNAME} \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
    --peerAddresses ${CORE_PEER_ADDRESS} \
    --tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE} \
    -C $CHANNEL_NAME -n ${CC_NAME} --isInit -c "${CC_INIT}"
