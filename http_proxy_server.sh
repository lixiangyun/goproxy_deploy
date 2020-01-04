#!/bin/bash
export PATH=$PATH:/usr/bin

goproxy_basic -addr 127.0.0.1:1000 &

tcpproxy -config config_server.yaml