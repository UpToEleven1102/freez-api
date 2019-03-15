#!/bin/bash

set -f
dev_server=$DEV_SERVER
environment=$ENV

echo "Deploy to dev server: ${dev_server}"
make push-code-gitlab
ssh -i ~/.ssh/id_rsa ubuntu@${dev_server} "cd ~/go/src/git.nextgencode.io/huyen.vu/freez-app-rest && sudo docker-compose up"
