FROM golang:latest
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux

RUN go get -u -v github.com/lixiangyun/goproxy_deploy
WORKDIR /gopath/src/github.com/lixiangyun/goproxy_deploy/server
RUN go build .

FROM ubuntu:xenial
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /usr/bin/
COPY --from=0 /gopath/src/github.com/lixiangyun/goproxy_deploy/server/server ./server

RUN chmod +x *

EXPOSE 8080

ENTRYPOINT ["server"]