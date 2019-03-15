#!/bin/bash

set -f
dev_server=$DEV_SERVER
environment=$ENV

echo "Deploy to dev server: ${dev_server}"
scp -r $CI_PROJECT_DIR ubuntu@35.162.158.187:/home/ubuntu/go/src/git.nextgencode.io/huyen.vu
ssh ubuntu@${dev_server} "cd ~/go/src/git.nextgencode.io/huyen.vu/freez-app-rest && sudo docker-compose kill && sudo docker-compose rm -v && sudo docker build -t freez-app-rest . && sudo docker-compose up"
