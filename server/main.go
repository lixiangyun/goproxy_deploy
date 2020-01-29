package main

import (
	"crypto/tls"
	"net"
	"flag"
	"log"
	proxy "github.com/lixiangyun/goproxy_deploy"
	"os"
	"sync"
)

var (
	Help bool
	ListenAddr  string

    ProxyStat *proxy.Stat
)

func init()  {
	flag.BoolVar(&Help,"help",false,"usage help.")
	flag.StringVar(&ListenAddr,"listen",":8080","listen address.")
}

func LocalConnect(addr string, localConn net.Conn)  {
	remoteConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	syncSem := new(sync.WaitGroup)
	syncSem.Add(2)
	go proxy.TcpChannel(ProxyStat, true, localConn, remoteConn, syncSem)
	go proxy.TcpChannel(ProxyStat, false, remoteConn, localConn, syncSem)
	syncSem.Wait()
}

func main()  {
	flag.Parse()
	if Help {
		flag.Usage()
		os.Exit(-1)
	}

	LocalProxy := proxy.ProxyBasic()
	log.Printf("start internal proxy basic %s\n", LocalProxy)

	listen,err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	tlscfg, err := proxy.TlsConfigServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	listen = tls.NewListener(listen,tlscfg)

	log.Printf("start public proxy basic %s\n", ListenAddr)

	ProxyStat = proxy.NewStat("stat")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("remote : %s\n",conn.RemoteAddr().String())
		go LocalConnect(LocalProxy,conn)
	}
}
