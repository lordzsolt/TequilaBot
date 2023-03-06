#!/usr/bin/env bash

set -xe

DESTINATION=lord@iosmith.com

GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux .

ssh $DESTINATION 'sudo systemctl stop saltbot.service'
# scp config.json $DESTINATION:/home/saltbot/config.json
scp bin/app-amd64-linux $DESTINATION:/home/saltbot/saltbot.bin
ssh $DESTINATION 'sudo systemctl start saltbot.service'