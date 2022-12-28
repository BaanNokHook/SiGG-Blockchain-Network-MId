#!/bin/bash -l

ROOT_PART=$PWD

NBTC_NETWORK_PATH=$ROOT_PART/network/nbtc
AIS_NETWORK_PATH=$ROOT_PART/network/ais
BBL_NETWORK_PATH=$ROOT_PART/network/bbl

# cd $ROOT_PART/chaincode && ./build.sh

cd $ROOT_PART/config && ./generate.sh

# cd $NBTC_NETWORK_PATH && sudo ./gen_tls_connection.sh &
# cd $AIS_NETWORK_PATH && sudo ./gen_tls_connection.sh &
# cd $BBL_NETWORK_PATH && sudo ./gen_tls_connection.sh &

# wait