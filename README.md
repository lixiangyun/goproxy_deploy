# HTTP二级代理(加密)

## 一级HTTP代理

```
docker run -d --restart=always --net=host linimbus/goproxy_deploy
```

## 二级HTTP代理
1、修改config_client.yaml的VIP地址；
2、启动本地二级代理；
```
http_proxy_client.bat
```
3、设置浏览器代理地址127.0.0.1:8080
