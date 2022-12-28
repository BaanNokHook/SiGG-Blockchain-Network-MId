
CHAINCODE_PATH=mobileidcc
# CHAINCODE_PATH=fabcar

# rm -rf mobileidcc.tar.gz &
# cd $CHAINCODE_PATH/src && rm -rf vendor &

# wait

cd $CHAINCODE_PATH/src && GO111MODULE=on go mod vendor &

wait

cd $CHAINCODE_PATH && tar cfz code.tar.gz src &

wait

cd $CHAINCODE_PATH && tar cfz mobileidcc.tar.gz code.tar.gz metadata.json &

wait

cd $CHAINCODE_PATH && mv mobileidcc.tar.gz ../ &

cd $CHAINCODE_PATH && rm -rf code.tar.gz &

wait