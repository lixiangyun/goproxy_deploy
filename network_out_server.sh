#!/bin/bash
cd /root/goproxy_deploy
nohup ./proxy bridge -p ":38000" -C ./proxy.crt -K ./proxy.key --log bridge.log &
nohup ./proxy server -r ":33389@:3389" -P "127.0.0.1:38000" -C ./proxy.crt -K ./proxy.key --log server3389.log &
nohup ./proxy server -r ":30080@:80" -P "127.0.0.1:38000" -C ./proxy.crt -K ./proxy.key --log server80.log &
nohup ./proxy server -r ":30021@:21" -P "127.0.0.1:38000" -C ./proxy.crt -K ./proxy.key --log server21.log &
exit 0