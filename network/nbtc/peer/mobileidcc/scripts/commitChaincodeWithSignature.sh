#!/bin/bash -l
CORE_PEER_ADDRESS_NBTC=$CORE_PEER_ADDRESS
CORE_PEER_TLS_ROOTCERT_FILE_NBTC=$CORE_PEER_TLS_ROOTCERT_FILE

CORE_PEER_ADDRESS_AIS=peer0.ais.mobileid.com:7051
CORE_PEER_TLS_ROOTCERT_FILE_AIS=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/ais.mobileid.com/peers/peer0.ais.mobileid.com/tls/ca.crt

CORE_PEER_ADDRESS_BBL=peer0.bbl.mobileid.com:7051
CORE_PEER_TLS_ROOTCERT_FILE_BBL=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bbl.mobileid.com/peers/peer0.bbl.mobileid.com/tls/ca.crt

peer lifecycle chaincode commit \
    -o ${ORDERER_ADDRESS} --ordererTLSHostnameOverride ${OVERIDE_ORDERER_HOSTNAME} \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
    --channelID $CHANNEL_NAME --name ${CC_NAME} \
    --peerAddresses ${CORE_PEER_ADDRESS_NBTC} \
    --tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE_NBTC} \
    --peerAddresses ${CORE_PEER_ADDRESS_AIS} \
    --tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE_AIS} \
    --peerAddresses ${CORE_PEER_ADDRESS_BBL} \
    --tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE_BBL} \
    --version ${CC_VERSION} --sequence ${CC_SEQUENCE} \
    --init-required \
    --signature-policy "${CC_POLICY}" --waitForEvent
