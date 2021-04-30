#!/usr/bin/env bash
set -xe

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
export GOPATH="$SCRIPT_DIR/vendor"

go build -o bin/bot .
