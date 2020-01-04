#!/bin/bash
export PATH=$PATH:/usr/bin

goproxy_basic -addr 127.0.0.1:1000 &
goproxy_basic -addr 127.0.0.1:1001 &
goproxy_basic -addr 127.0.0.1:1002 &
goproxy_basic -addr 127.0.0.1:1003 &
goproxy_basic -addr 127.0.0.1:1004 &

tcpproxy -config config_server.yaml