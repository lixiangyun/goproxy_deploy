### http代理服务端镜像，默认监听：8080 端口

- docker build -t linimbus/goproxy_docker -f Dockerfile_goproxy .
- docker run -d -p 8080:8080 --restart=always linimbus/goproxy_docker

### 在客户端设置HTTP代理环境变量

- http_proxy=http://xxx.xxx.xxx.xxx:8080
- https_proxy=http://xxx.xxx.xxx.xxx:8080
