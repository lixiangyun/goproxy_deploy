#!/bin/bash
cd /root/goproxy_deploy
killall goproxy_basic
killall tcpproxy
nohup ./goproxy_basic -addr 127.0.0.1:1000 &
nohup ./goproxy_basic -addr 127.0.0.1:1001 &
nohup ./goproxy_basic -addr 127.0.0.1:1002 &
nohup ./goproxy_basic -addr 127.0.0.1:1003 &
nohup ./goproxy_basic -addr 127.0.0.1:1004 &
nohup ./tcpproxy -config config_server.yaml &
exit 0
