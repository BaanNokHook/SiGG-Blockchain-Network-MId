export FABRIC_CFG_PATH=$PWD

#====================== Genesis Block ===================================//
../bin/configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID orderer-system-channel

#====================== Channel Config ===================================//
../bin/configtxgen -profile MobileIDPeerChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mobileid
../bin/configtxgen -profile MobileIDPeerChannel -outputAnchorPeersUpdate ./channel-artifacts/NBTCMSPanchors.tx -channelID mobileid -asOrg NBTC
../bin/configtxgen -profile MobileIDPeerChannel -outputAnchorPeersUpdate ./channel-artifacts/AISMSPanchors.tx -channelID mobileid -asOrg AIS
../bin/configtxgen -profile MobileIDPeerChannel -outputAnchorPeersUpdate ./channel-artifacts/BBLMSPanchors.tx -channelID mobileid -asOrg BBL
