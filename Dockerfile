FROM ubuntu:xenial
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /usr/bin/

ADD ./crt/client.crt ./crt/client.crt
ADD ./crt/client.pem ./crt/client.pem
ADD ./crt/server.crt ./crt/server.crt
ADD ./crt/server.pem ./crt/server.pem

ADD config_server.yaml config_server.yaml
ADD http_proxy_server.sh http_proxy_server.sh

ADD goproxy_basic goproxy_basic
ADD tcpproxy tcpproxy

RUN chmod +x *

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["http_proxy_server.sh"]