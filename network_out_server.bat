start proxy.exe bridge -p ":38000" -C ./proxy.crt -K ./proxy.key
start proxy.exe server -r ":33389@:3389" -P "127.0.0.1:38000" -C ./proxy.crt -K ./proxy.key
start proxy.exe server -r ":30080@:80" -P "127.0.0.1:38000" -C ./proxy.crt -K ./proxy.key
start proxy.exe server -r ":30021@:21" -P "127.0.0.1:38000" -C ./proxy.crt -K ./proxy.key
pause