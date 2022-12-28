#!/bin/bash -l

peer chaincode upgrade --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA \
  -n ${CC_NAME} -v ${CC_VERSION} \
  -c "${CC_INIT}" \
  -p ${CC_SRC_PATH} -C ${CHANNEL_NAME} \
  -o ${ORDERER_ADDRESS} \
  -P "${CC_POLICY}"