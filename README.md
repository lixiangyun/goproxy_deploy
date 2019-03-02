# HTTP二级代理(加密)

## 一级HTTP代理(VPS,IP:22.22.22.22)
```
./proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key
```

## 二级HTTP代理(本地Linux)
```
./proxy http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key
```
那么访问本地的8080端口就是访问VPS上面的代理端口38080.

## 二级HTTP代理(本地windows)
```
./proxy.exe http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key
```

然后设置你的windos系统中，需要通过代理上网的程序的代理为http模式，地址为：127.0.0.1，端口为：8080,程序即可通过加密通道通过vps上网。

## 安装到linux后台服务
```
cp goproxy.service /lib/systemd/system/goproxy.service 
systemctl daemon-reload
service goproxy start
```

## 内网穿透

### 公司机器A提供了web服务80端口，有VPS一个,公网IP:22.22.22.22

需求:
在家里能够通过访问VPS的28080端口访问到公司机器A的80端口

步骤:

- 在vps上执行

    ```
    ./proxy bridge -p ":33080" -C proxy.crt -K proxy.key
    ./proxy server -r ":28080@:80" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key
    ```

- 在公司机器A上面执行
        
    ```
    ./proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key
    ```
