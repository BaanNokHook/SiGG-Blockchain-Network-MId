#!/bin/bash -l

peer chaincode invoke -o ${ORDERER_ADDRESS} --tls true --cafile ${ORDERER_CA} \
-C ${CHANNEL_NAME} -n ${CC_NAME} \
--waitForEvent \
--peerAddresses ${CORE_PEER_ADDRESS} \
--tlsRootCertFiles ${CORE_PEER_TLS_ROOTCERT_FILE} \
-c "${CC_INVOCATION}"
