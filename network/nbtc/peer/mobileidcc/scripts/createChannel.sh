# CHANNEL_NAME=mobileid
# ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/mobileid.com/orderers/orderer.mobileid.com/msp/tlscacerts/tlsca.mobileid.com-cert.pem
peer channel create -o $ORDERER_ADDRESS -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile $ORDERER_CA
