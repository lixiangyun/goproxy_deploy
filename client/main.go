package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	BlockAddressPools map[string]string
	BlockLock sync.RWMutex
)

var (
	Help bool
	RemoteProxy string
	ListenAddr  string
)

func init()  {
	BlockAddressPools = make(map[string]string,0)

	flag.BoolVar(&Help,"help",false,"usage help.")
	flag.StringVar(&RemoteProxy,"proxy","www.proxy.com:8080","remote proxy address.")
	flag.StringVar(&ListenAddr,"listen",":8080","listen address.")
}

func IsBlock(address string) bool {
	BlockLock.RLock()
	_,flag := BlockAddressPools[address]
	BlockLock.RUnlock()
	if flag == true {
		return true
	}
	
	flag = retryConn(address)
	if flag == true {
		return false
	}
	
	BlockLock.Lock()
	BlockAddressPools[address] = address
	BlockLock.Unlock()

	log.Println("block address : " + address)

	return true
}

func retryConn(address string) bool {
	 conn, err := net.DialTimeout("tcp",address,1*time.Second)
	 if err != nil {
	 	log.Println(err.Error())
	 	return false
	 }
	 conn.Close()
	 return true
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
		log.Printf("host : %s\n", addr)
		return buffer[:begin],addr
	}
}

func writeFull(conn net.Conn, buf []byte) error {
	totallen := len(buf)
	sendcnt := 0

	for {
		cnt, err := conn.Write(buf[sendcnt:])
		if err != nil {
			return err
		}
		if cnt+sendcnt >= totallen {
			return nil
		}
		sendcnt += cnt
	}
}

// tcp通道互通
func tcpChannel(s *Stat, up bool,localConn net.Conn, remoteConn net.Conn, wait *sync.WaitGroup) {
	defer wait.Done()
	defer localConn.Close()
	defer remoteConn.Close()

	buf := make([]byte, 65535)
	for {
		cnt, err := localConn.Read(buf[0:])
		if err != nil {
			if cnt != 0 {
				writeFull(remoteConn, buf[0:cnt])
			}
			break
		}
		if up {
			s.Add(cnt,0)
		}else {
			s.Add(0,cnt)
		}
		err = writeFull(remoteConn, buf[0:cnt])
		if err != nil {
			break
		}
	}
}

func RemoteConnect(addr string, buffer []byte, localConn net.Conn)  {
	remoteConn, err := net.Dial("tcp",addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	tlsCfg,err := TlsConfigClient(nil,addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	remoteConn = tls.Client(remoteConn, tlsCfg)

	err = writeFull(remoteConn,buffer)
	if err != nil {
		log.Println(err.Error())
		remoteConn.Close()
		return
	}

	syncSem := new(sync.WaitGroup)
	syncSem.Add(2)
	go tcpChannel(remoteProxyStat, true, localConn, remoteConn, syncSem)
	go tcpChannel(remoteProxyStat, false, remoteConn, localConn, syncSem)
	syncSem.Wait()
}

func LocalConnect(addr string, buffer []byte, localConn net.Conn)  {
	remoteConn, err := net.Dial("tcp",addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = writeFull(remoteConn,buffer)
	if err != nil {
		log.Println(err.Error())
		remoteConn.Close()
		return
	}

	syncSem := new(sync.WaitGroup)
	syncSem.Add(2)
	go tcpChannel(localProxyStat, true, localConn, remoteConn, syncSem)
	go tcpChannel(localProxyStat, false, remoteConn, localConn, syncSem)
	syncSem.Wait()
}

var localProxyStat *Stat
var remoteProxyStat *Stat

func main()  {
	flag.Parse()
	if Help {
		flag.Usage()
		os.Exit(-1)
	}

	LocalProxy := proxyBasic()
	log.Printf("start internal proxy basic %s\n", LocalProxy)

	listen,err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("start public proxy basic %s\n", ListenAddr)

	localProxyStat = NewStat("local")
	remoteProxyStat = NewStat("remote")

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
		if true == IsBlock(addr) {
			go RemoteConnect(RemoteProxy,buffer,conn)
		}else {
			go LocalConnect(LocalProxy,buffer,conn)
		}
	}
}