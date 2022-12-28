#!/bin/bash -l

NEW_ORG=$1
NEW_MSP=$2

echo
echo "========= Config transaction to update channel config ===== "
echo

echo "Signing config transaction"
echo

peer channel signconfigtx -f channel_update_in_envelope.pb

echo
echo "========= Submitting channel update transaction ========= "
echo
set -x
peer channel update -f channel_update_in_envelope.pb -c ${CHANNEL_NAME} -o ${ORDERER_ADDRESS} --tls --cafile ${ORDERER_CA}
set +x

echo
echo "========= Config transaction to update channel config submitted! =========== "
echo

exit 0
