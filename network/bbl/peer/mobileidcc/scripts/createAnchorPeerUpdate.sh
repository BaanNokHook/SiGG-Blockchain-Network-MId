#!/bin/bash -l

. ./scripts/fetchChannelConfig.sh
. ./scripts/createConfigUpdate.sh

# Fetch the config for the channel, writing it to config.json
fetchChannelConfig ${CHANNEL_NAME} config.json

# Modify the configuration to add new org artifacts
ANCHOR_PEER_ADDRESS=$(echo ${CORE_PEER_ADDRESS} | sed 's/:.*//g' --)
ANCHOR_PEER_PORT=$(echo ${CORE_PEER_ADDRESS} | sed 's/.*://g' --)
jq '.channel_group.groups.Application.groups.'${CAP_ORGNAME}'.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "'${ANCHOR_PEER_ADDRESS}'","port": '${ANCHOR_PEER_PORT}'}]},"version": "0"}}' config.json > modified_anchor_config.json

# Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to channel_update_in_envelope.pb
createConfigUpdate ${CHANNEL_NAME} config.json modified_anchor_config.json channel_update_in_envelope.pb
