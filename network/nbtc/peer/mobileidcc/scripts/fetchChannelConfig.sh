#!/bin/bash -l

# fetchChannelConfig <channel_id> <output_json>
# Writes the current channel config for a given channel to a JSON file
fetchChannelConfig() {
  CHANNEL=$1
  OUTPUT=$2

  echo
  echo "========= Fetch channel config from orderer =========== "
  echo

  echo "Fetching the most recent configuration block for the channel"
  set -x
  peer channel fetch config config_block.pb -o ${ORDERER_ADDRESS} -c ${CHANNEL_NAME} --tls --cafile ${ORDERER_CA}
  set +x

  echo "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  set +x
}


## Fetch the config for the channel, writing it to config.json
# fetchChannelConfig ${CHANNEL_NAME} config.json

