#!/usr/bin/env bash

# Mount EFS
sudo yum install -y amazon-efs-utils
sudo mkdir /efs
sudo chown webapp /efs
sudo mount -t efs -o tls fs-41f8231a:/ /efs

set -xe

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
export GOPATH="$SCRIPT_DIR/vendor"

go build -o bin/bot .
