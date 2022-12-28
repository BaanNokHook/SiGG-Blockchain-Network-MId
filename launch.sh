#!/bin/bash -l

ROOT_PART=$PWD

NBTC_NETWORK_PATH=$ROOT_PART/network/nbtc
AIS_NETWORK_PATH=$ROOT_PART/network/ais
BBL_NETWORK_PATH=$ROOT_PART/network/bbl

# cd $NBTC_NETWORK_PATH && ./package_chaincode.sh &
# cd $AIS_NETWORK_PATH && ./package_chaincode.sh &
# cd $BBL_NETWORK_PATH && ./package_chaincode.sh &

# wait

cd $NBTC_NETWORK_PATH && ./create_channel.sh &

wait

cd $NBTC_NETWORK_PATH && ./joining_channel.sh &
cd $AIS_NETWORK_PATH && ./joining_channel.sh &
cd $BBL_NETWORK_PATH && ./joining_channel.sh &

wait

cd $NBTC_NETWORK_PATH && ./install_chaincode.sh &
cd $AIS_NETWORK_PATH && ./install_chaincode.sh &
cd $BBL_NETWORK_PATH && ./install_chaincode.sh &

wait

# cd $NBTC_NETWORK_PATH && ./approve_chaincode.sh &
# cd $AIS_NETWORK_PATH && ./approve_chaincode.sh &
# cd $BBL_NETWORK_PATH && ./approve_chaincode.sh &

cd $NBTC_NETWORK_PATH && ./approve_chaincode_with_signature.sh &
cd $AIS_NETWORK_PATH && ./approve_chaincode_with_signature.sh &
cd $BBL_NETWORK_PATH && ./approve_chaincode_with_signature.sh &

wait

# cd $NBTC_NETWORK_PATH && ./check_commit_readiness.sh &

cd $NBTC_NETWORK_PATH && ./check_commit_readiness_with_signature.sh &

wait

# cd $NBTC_NETWORK_PATH && ./commit_chaincode.sh &

cd $NBTC_NETWORK_PATH && ./commit_chaincode_with_signature.sh &

wait

cd $NBTC_NETWORK_PATH && ./query_committed.sh &

wait

cd $NBTC_NETWORK_PATH && ./init_ledger.sh &

wait

sleep 5

cd $NBTC_NETWORK_PATH && ./invoke_chaincode.sh &
cd $AIS_NETWORK_PATH && ./invoke_chaincode.sh &
cd $BBL_NETWORK_PATH && ./invoke_chaincode.sh &

wait

cd $NBTC_NETWORK_PATH && ./query_chaincode.sh
# cd $AIS_NETWORK_PATH && ./query_chaincode.sh
# cd $BBL_NETWORK_PATH && ./query_chaincode.sh

wait