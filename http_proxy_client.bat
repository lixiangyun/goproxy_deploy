#start proxy.exe http -t tcp -p ":8080" -T tcp -P "127.0.0.1:8000"
tcpproxy.exe -config config_client.yaml
pause
#start vchannel.exe -config config_channel.yaml -mode client