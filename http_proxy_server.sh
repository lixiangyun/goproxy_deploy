#!/bin/bash
cd /root/goproxy_deploy
killall goproxy_basic
killall vchannel
nohup ./goproxy_basic -addr :8080 &
nohup ./goproxy_basic -addr :8081 &
nohup ./goproxy_basic -addr :8082 &
nohup ./goproxy_basic -addr :8083 &
nohup ./goproxy_basic -addr :8084 &
nohup ./goproxy_basic -addr :8085 &
nohup ./goproxy_basic -addr :8086 &
nohup ./goproxy_basic -addr :8087 &
nohup ./goproxy_basic -addr :8088 &
nohup ./goproxy_basic -addr :8089 &
nohup ./vchannel -config config_channel.yaml -mode server &
exit 0
