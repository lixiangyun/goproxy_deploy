#!/bin/bash
apt update
apt install docker.io -y
docker run -d --restart=always --net=host linimbus/goproxy_deploy 