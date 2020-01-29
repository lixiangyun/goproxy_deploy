package proxy

import (
	"log"
	"net"
	"sync"
)

func WriteFull(conn net.Conn, buf []byte) error {
	totalLen := len(buf)
	sendCnt := 0

	for {
		cnt, err := conn.Write(buf[sendCnt:])
		if err != nil {
			return err
		}
		if cnt+sendCnt >= totalLen {
			return nil
		}
		sendCnt += cnt
	}
}

// tcp通道互通
func TcpChannel(s *Stat, up bool,localConn net.Conn, remoteConn net.Conn, wait *sync.WaitGroup) {
	defer wait.Done()
	defer localConn.Close()
	defer remoteConn.Close()

	buf := make([]byte, 65535)
	for {
		cnt, err := localConn.Read(buf[0:])
		if err != nil {
			if cnt != 0 {
				WriteFull(remoteConn, buf[0:cnt])
			}
			log.Println(err.Error())
			break
		}
		if up {
			s.Add(cnt,0)
		}else {
			s.Add(0,cnt)
		}
		err = WriteFull(remoteConn, buf[0:cnt])
		if err != nil {
			log.Println(err.Error())
			break
		}
	}
}