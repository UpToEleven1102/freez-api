#!/bin/bash

set -f
dev_server=$DEV_SERVER
environment=$ENV

echo "Deploy to dev server: ${dev_server}"
echo "Building exe file"
ssh ubuntu@${dev_server} "cd ~/go/src/git.nextgencode.io/huyen.vu/freez-app-rest && git pull origin develop && echo PATH=$PATH:/usr/local/go/bin && make build && sudo docker-compose kill && sudo docker-compose rm -v && sudo docker build -t freez-app-rest . && sudo docker-compose up -d"
