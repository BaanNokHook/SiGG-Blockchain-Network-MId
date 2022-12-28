#!/bin/bash -l

ROOT_PART=$PWD

NBTC_NETWORK_PATH=$ROOT_PART/network/nbtc
AIS_NETWORK_PATH=$ROOT_PART/network/ais
BBL_NETWORK_PATH=$ROOT_PART/network/bbl

# ./generate.sh

# wait

cd $NBTC_NETWORK_PATH && ./deploy_orderer_services.sh &
# cd $NBTC_NETWORK_PATH && ./deploy_ca_services.sh &
cd $AIS_NETWORK_PATH && ./deploy_orderer_services.sh &
cd $AIS_NETWORK_PATH && ./deploy_ca_services.sh &
cd $BBL_NETWORK_PATH && ./deploy_orderer_services.sh &
# cd $BBL_NETWORK_PATH && ./deploy_ca_services.sh &

wait

# cd $AIS_NETWORK_PATH && ./deploy_api_services.sh &
cd $NBTC_NETWORK_PATH && ./deploy_peer_services.sh &
cd $AIS_NETWORK_PATH && ./deploy_peer_services.sh &
cd $BBL_NETWORK_PATH && ./deploy_peer_services.sh &

wait