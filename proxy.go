package proxy

import (
	"log"
	"fmt"
	"github.com/lixiangyun/goproxy"
	"net"
	"net/http"
)

func checkAddress() string {
	beginPort := 8090
	for {
		address := fmt.Sprintf("127.0.0.1:%d",beginPort)
		listen,err := net.Listen("tcp",address)
		if err != nil {
			beginPort++
			continue
		}
		listen.Close()
		return address
	}
}

func ProxyBasic() string {
	address := checkAddress()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	go func() {
		log.Fatal(http.ListenAndServe(address, proxy))
	}()
	return address
}

