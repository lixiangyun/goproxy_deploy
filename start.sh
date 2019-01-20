cd /root/goproxy_deploy
./proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key --forever --log proxy.log --daemon