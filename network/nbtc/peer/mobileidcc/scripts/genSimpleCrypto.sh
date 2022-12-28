#!/bin/bash -l

FABRIC_CFG_PATH=$PWD
cryptogen generate --config=$PWD/test-org-certs.yaml
