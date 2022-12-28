#!/bin/bash -l

# fetchChannelConfig <channel_id> <output_json>
# Writes the current channel config for a given channel to a JSON file
fetchChannelConfig() {
  CHANNEL=$1
  OUTPUT=$2

  echo "Fetching the most recent configuration block for the channel"
  set -x
  peer channel fetch config config_block.pb -o ${ORDERER_ADDRESS} -c ${CHANNEL_NAME} --tls --cafile ${ORDERER_CA}
  set +x

  echo "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  set +x
}

# createConfigUpdate <channel_id> <original_config.json> <modified_config.json> <output.pb>
# Takes an original and modified config, and produces the config update tx
# which transitions between the two
createConfigUpdate() {
  CHANNEL_NAME=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4

  set -x
  configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb
  configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb
  configtxlator compute_update --channel_id "${CHANNEL_NAME}" --original original_config.pb --updated modified_config.pb >config_update.pb
  configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'${CHANNEL_NAME}'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json
  configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}"
  set +x
}

echo
echo "========= Creating config transaction to update channel config =========== "
echo

# Fetch the config for the channel, writing it to config.json
fetchChannelConfig ${CHANNEL_NAME} config.json

# Modify the configuration to change batch timeout
set -x
JQ_QUERY=".channel_group.groups.Orderer.values.BatchTimeout.value.timeout"
jq "${JQ_QUERY} = \"1s\"" config.json > modified_config.json
jq "${JQ_QUERY}" modified_config.json
set +x

# Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to channel_update_in_envelope.pb
createConfigUpdate ${CHANNEL_NAME} config.json modified_config.json channel_update_in_envelope.pb

echo
echo "========= Config transaction to update channel config ===== "
echo

echo "Signing config transaction"
echo

# signing as Orderer Admin
CORE_PEER_LOCALMSPID="OrdererMSP" \
  CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/mobileid.com/orderers/orderer3.mobileid.com/tls/ca.crt \
  CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/mobileid.com/users/Admin@mobileid.com/msp \
  CORE_PEER_ADDRESS=orderer3.mobileid.com:7050 \
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
