#!/bin/bash -l

peer lifecycle chaincode checkcommitreadiness \
    --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} \
    --output json --init-required \
    --signature-policy "${CC_POLICY}"