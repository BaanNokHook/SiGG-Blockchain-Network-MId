#!/bin/bash -l

# Generate channel configuration transaction
function generateChannelArtifacts() {
  NEW_ORG=$1
  NEW_MSP=$2

  which configtxgen
  if [ "$?" -ne 0 ]; then
    echo "configtxgen tool not found. exiting"
    exit 1
  fi
  echo "##########################################################"
  echo "#########  Generating PeerOrg config material ###########"
  echo "##########################################################"
  export FABRIC_CFG_PATH=$PWD
  set -x
  configtxgen -printOrg $NEW_ORG > $PWD/channel-artifacts/$NEW_ORG.json
  res=$?

  set +x
  if [ $res -ne 0 ]; then
   echo "Failed to generate PeerOrg config material..."
   exit 1
  fi
  echo
}

# generateChannelArtifacts $NEW_ORG $NEW_MSP
