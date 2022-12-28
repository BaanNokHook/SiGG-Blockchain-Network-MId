#!/bin/bash -l

NEW_ORG=$1
NEW_MSP=$2

echo
echo "========= Creating config transaction to update channel config =========== "
echo
. ./scripts/fetchChannelConfig.sh
. ./scripts/generateChannelArtifacts.sh
. ./scripts/modifyConfig.sh
. ./scripts/createConfigUpdate.sh

# Fetch the config for the channel, writing it to config.json
fetchChannelConfig ${CHANNEL_NAME} config.json

# generate config that contains new org artifacts
generateChannelArtifacts $NEW_ORG $NEW_MSP

# Modify the configuration to add new org artifacts
modifyConfig $NEW_ORG config.json "$PWD/channel-artifacts/$NEW_ORG.json" modified_config.json

# Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to channel_update_in_envelope.pb
createConfigUpdate ${CHANNEL_NAME} config.json modified_config.json channel_update_in_envelope.pb

