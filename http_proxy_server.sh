#!/bin/bash
cd /root/goproxy_deploy
./proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon
./proxy http -t tls -p ":38081" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon
./proxy http -t tls -p ":38082" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon
./proxy http -t tls -p ":38083" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon
./proxy http -t tls -p ":38084" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon
./proxy http -t tls -p ":38085" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon
exit 0
