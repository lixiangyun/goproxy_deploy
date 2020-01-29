package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	proxy "github.com/lixiangyun/goproxy_deploy"
)

var (
	BlockAddressList map[string]bool
	BlockLock sync.RWMutex
)

var (
	Help bool
	BlockCacheFile string
	RemoteProxy string
	ListenAddr  string
	RetryTimeOut int
)

func init()  {
	BlockAddressList = make(map[string]bool,0)

	flag.BoolVar(&Help,"help",false,"usage help.")
	flag.IntVar(&RetryTimeOut,"retry",5,"retry connect timeout, default 5 seconds.")
	flag.StringVar(&RemoteProxy,"proxy","www.proxy.com:8080","remote proxy address.")
	flag.StringVar(&ListenAddr,"listen",":8080","listen address.")
	flag.StringVar(&BlockCacheFile,"db","black.db","black/white address list.")
}

func parseBlack(line string) (string,bool) {
	hosts := strings.Split(line,"\t")
	if len(hosts) != 2 {
		return "", false
	}
	if 0 == strings.Compare(hosts[1],"true") {
		return hosts[0], true
	}
	return hosts[0], false
}

func loadBlackCacheFile()  {
	body, err := ioutil.ReadFile(BlockCacheFile)
	if err != nil {
		log.Println(err.Error())
		return
	}
	lines := strings.Split(string(body),"\n")
	for _,v := range lines {
		address,blackFlag := parseBlack(v)
		if address == "" {
			continue
		}
		BlockAddressList[address] = blackFlag
	}
}

func saveBlackCacheFile(address string, black bool)  {
	BlockLock.Lock()
	defer BlockLock.Unlock()

	file,err := os.OpenFile(BlockCacheFile,os.O_APPEND|os.O_WRONLY,0644)
	if err != nil {
		file, err = os.Create(BlockCacheFile)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	defer file.Close()

	if black {
		fmt.Fprintf(file,"%s\t%s\n",address,"true")
	}else {
		fmt.Fprintf(file,"%s\t%s\n",address,"false")
	}
}

func IsBlack(address string) bool {

	BlockLock.RLock()
	black,flag := BlockAddressList[address]
	BlockLock.RUnlock()
	if flag == true {
		return black
	}

	black = retryConn(address)

	BlockLock.Lock()
	BlockAddressList[address] = black
	BlockLock.Unlock()

	saveBlackCacheFile(address,black)

	return black
}

func retryConn(address string) bool {
	 conn, err := net.DialTimeout("tcp",address,time.Duration(RetryTimeOut)*time.Second)
	 if err != nil {
	 	log.Println(err.Error())
	 	return true
	 }
	 conn.Close()
	 return false
}

func parseAddress(buffer []byte) string {
	line := strings.Split(string(buffer),"\r\n")
	for _,v := range line {
		if strings.HasPrefix(v,"Host: ") {
			if strings.Index(v[6:],":") == -1 {
				return v[6:] + ":80"
			}else {
				return v[6:]
			}
		}
	}
	return ""
}

func loadAddress(conn net.Conn) ([]byte,string) {
	var buffer [4096]byte
	var begin int
	for {
		cnt,err := conn.Read(buffer[begin:])
		if err != nil {
			return buffer[:begin],""
		}
		begin += cnt

		addr := parseAddress(buffer[:begin])
		if addr == "" || len(addr) == 0 {
			continue
		}
		return buffer[:begin],addr
	}
}


func RemoteConnect(addr string, buffer []byte, localConn net.Conn)  {
	remoteConn, err := net.Dial("tcp",addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	tlsCfg,err := proxy.TlsConfigClient(addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	remoteConn = tls.Client(remoteConn, tlsCfg)

	err = proxy.WriteFull(remoteConn,buffer)
	if err != nil {
		log.Println(err.Error())
		remoteConn.Close()
		return
	}

	syncSem := new(sync.WaitGroup)
	syncSem.Add(2)
	go proxy.TcpChannel(remoteProxyStat, true, localConn, remoteConn, syncSem)
	go proxy.TcpChannel(remoteProxyStat, false, remoteConn, localConn, syncSem)
	syncSem.Wait()
}

func LocalConnect(addr string, buffer []byte, localConn net.Conn)  {
	remoteConn, err := net.Dial("tcp",addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = proxy.WriteFull(remoteConn,buffer)
	if err != nil {
		log.Println(err.Error())
		remoteConn.Close()
		return
	}

	syncSem := new(sync.WaitGroup)
	syncSem.Add(2)
	go proxy.TcpChannel(localProxyStat, true, localConn, remoteConn, syncSem)
	go proxy.TcpChannel(localProxyStat, false, remoteConn, localConn, syncSem)
	syncSem.Wait()
}

var localProxyStat *proxy.Stat
var remoteProxyStat *proxy.Stat

func main()  {
	flag.Parse()
	if Help {
		flag.Usage()
		os.Exit(-1)
	}

	loadBlackCacheFile()

	LocalProxy := proxy.ProxyBasic()
	log.Printf("start internal proxy basic %s\n", LocalProxy)

	listen,err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("start public proxy basic %s\n", ListenAddr)

	localProxyStat = proxy.NewStat("local")
	remoteProxyStat = proxy.NewStat("remote")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		buffer, addr := loadAddress(conn)
		if addr == "" {
			conn.Close()
			continue
		}
		if true == IsBlack(addr) {
			log.Printf("host : %s black\n", addr)
			go RemoteConnect(RemoteProxy,buffer,conn)
		}else {
			log.Printf("host : %s not black\n", addr)
			go LocalConnect(LocalProxy,buffer,conn)
		}
	}
}