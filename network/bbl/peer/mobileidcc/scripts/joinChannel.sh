#!/bin/bash -l

MAX_RETRY=10
DELAY=3

# Fetch 0 block for channel joining
fetchChannel(){
  set -x
  peer channel fetch 0 $CHANNEL_NAME.block -o ${ORDERER_ADDRESS} -c $CHANNEL_NAME --tls --cafile $ORDERER_CA >&info.txt
  res=$?
  set +x
}

## Sometimes Join takes time hence RETRY at least 5 times
joinChannel(){
  set -x
  peer channel join -b $CHANNEL_NAME.block >&info.txt
  res=$?
  set +x
}

fetchChannelWithRetry(){
  for (( i=1; i <= $MAX_RETRY; i++ ))
  do
    res=1 # init exit code
    echo "$i attempts to fetch channel with delay ${DELAY} sec."
    fetchChannel &>log.txt
    cat info.txt

    # Assert channel fetching
    if [ $res -eq 0 ]
    then
      echo "Fetch channel ${CHANNEL_NAME} successful."
      break
    elif [ $res -ne 0 ] && [ $i -le $MAX_RETRY ]
    then
      echo "WARNING: Cannot fetch channel ${CHANNEL_NAME}."
    fi
    if [ $res -ne 0 ] && [ $i -ge $MAX_RETRY ]
    then
      echo "ERROR: Failed to fetch channel ${CHANNEL_NAME} after retry ${MAX_RETRY} attempts."
      exit 1
    fi

    sleep ${DELAY}
  done
}

joinChannelWithRetry(){
  for (( i=1; i <= $MAX_RETRY; i++ ))
  do
    res=1 # init exit code
    echo "$i attempts to join channel with delay ${DELAY} sec."
    joinChannel &>log.txt
    cat info.txt

    # Assert channel joining
    if [ $res -eq 0 ]
    then
      echo "Join channel ${CHANNEL_NAME} successful."
      break
    elif [ $res -ne 0 ] && [ $i -le $MAX_RETRY ]
    then
      echo "WARNING: Cannot join channel ${CHANNEL_NAME}."
    fi
    if [ $res -ne 0 ] && [ $i -ge $MAX_RETRY ]
    then
      echo "ERROR: Failed to join channel ${CHANNEL_NAME} after retry ${MAX_RETRY} attempts."
      exit 1
    fi

    sleep ${DELAY}
  done
}

echo "========= Start channel joining phase ========="

fetchChannelWithRetry

joinChannelWithRetry

echo "========= Done ========="

exit 0
