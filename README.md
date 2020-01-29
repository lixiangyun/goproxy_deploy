# HTTP二级代理(TLS1.2加密通道)

本地二级代理，默认会尝试连接，如果本地访问失败。才会走远程代理。默认采用TLS1.2加密传输；
本地会生成black.db文件，用于缓存地址的可访问性。

## 一级HTTP代理
```
docker run -d --restart=always --net=host linimbus/goproxy_deploy
```

## 二级HTTP代理
1、启动本地二级代理；
```
client.exe -proxy {IP}:8080
```

2、设置浏览器代理地址127.0.0.1:8080
