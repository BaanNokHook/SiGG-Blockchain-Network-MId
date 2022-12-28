#!/bin/bash -l

DEFAULT_CAP_ORGNAME="$(echo "${CORE_PEER_LOCALMSPID}" | sed "s/MSP//" --)"
CAP_ORGNAME=${CAP_ORGNAME:=${DEFAULT_CAP_ORGNAME}}
. ./scripts/createAnchorPeerUpdate.sh ${CAP_ORGNAME:=${DEFAULT_CAP_ORGNAME}}
peer channel update -f channel_update_in_envelope.pb -o $ORDERER_ADDRESS -c $CHANNEL_NAME  --tls --cafile $ORDERER_CA