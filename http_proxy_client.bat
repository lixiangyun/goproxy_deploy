#start proxy.exe http -t tcp -p ":8080" -T tcp -P "127.0.0.1:8000"
start tcpproxy.exe -config config_client.yaml

#start vchannel.exe -config config_channel.yaml -mode client