cd /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode/mobileidcc/src/ && GO111MODULE=on go mod vendor &

wait

peer lifecycle chaincode package mobileidcc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode/mobileidcc/src/ --lang golang --label mobileidcc