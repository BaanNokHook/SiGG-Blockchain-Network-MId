#!/bin/bash -l

if [[ "$PATH" != *"$PWD/bin"* ]]; 
then
    echo "fabric ca client is not in PATH. appending ${PWD}/bin to PATH"
    export PATH=$PATH:$PWD/bin
    echo "export PATH=$PATH:$PWD/bin" >> ~/.bashrc
    echo "export PATH=$PATH:$PWD/bin" >> ~/.bash_profile
fi
if [[ -f $PWD/bin/fabric-ca-client ]]; 
then
    echo "fabric-ca-client already exists."
else
    echo "fabric-ca-client not found. installing fabric-ca-client binary ..."
    (
        cd ${PROJECT_DIR}
        curl -sSL https://bit.ly/2ysbOFE | bash -s -- -d -s 2.1.1 1.4.7 0.4.20
    )
fi
echo "using fabric-ca-client from" $(which fabric-ca-client)
