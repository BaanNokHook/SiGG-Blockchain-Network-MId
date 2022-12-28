#!/bin/bash -l

peer chaincode query -C ${CHANNEL_NAME} -n ${CC_NAME} -c "${CC_QUERY}"
