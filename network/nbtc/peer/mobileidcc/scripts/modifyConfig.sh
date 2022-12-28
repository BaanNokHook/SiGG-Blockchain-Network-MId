#!/bin/bash -l

# Modify the configuration to update
function modifyConfig() {
  TARGET_MSP=$1
  ORIGIN_CONFIG=$2
  NEW_CONFIG=$3
  UPDATED_CONFIG=$4

  set -x
  JQ_QUERY=".[0] * {\"channel_group\":{\"groups\":{\"Application\":{\"groups\": {\"${TARGET_MSP}\":.[1]}}}}}"
  jq -s "${JQ_QUERY}" ${ORIGIN_CONFIG} ${NEW_CONFIG} >${UPDATED_CONFIG}
  set +x
}
