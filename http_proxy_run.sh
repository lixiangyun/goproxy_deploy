#!/bin/bash
cd /root/goproxy_deploy
nohup ./proxy http -t tcp -p ":8080" -T tls -P "aws.lixiangyun.top:38080" -C proxy.crt -K proxy.key &
exit 0
