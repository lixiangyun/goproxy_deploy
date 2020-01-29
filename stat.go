package proxy

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type Stat struct {
	up uint64
	down uint64
	prefix string
}

func NewStat(prefix string) *Stat {
	s := new(Stat)
	s.prefix = prefix
	go s.display()
	return s
}

func (s *Stat)Add(up int, down int) {
	if up > 0 {
		atomic.AddUint64(&s.up, uint64(up))
	}
	if down > 0 {
		atomic.AddUint64(&s.down, uint64(down))
	}
}

func (s *Stat)display() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C
			log.Printf("%s up:%s down:%s\n",
				s.prefix, calcUnit(s.up), calcUnit(s.down))
		}
	}()
}

func calcUnit(cnt uint64) string {
	if cnt < 1024 {
		return fmt.Sprintf("%d", cnt)
	} else if cnt < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float32(cnt)/1024)
	} else if cnt < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float32(cnt)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float32(cnt)/(1024*1024*1024))
	}
}
