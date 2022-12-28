export FABRIC_CFG_PATH=$PWD

./cleanArtifacts.sh

./cleanCrypto.sh

../bin/cryptogen extend --config=./crypto-config.yaml

./createProfile.sh
