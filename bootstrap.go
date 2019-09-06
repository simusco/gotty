package gotty

import (
	"log"
	"net"
	"sync"
	"time"
)

type ChannelInitializer func(*Channel)

type Bootstrap struct {
	Addr            string
	MaxConnNum      int
	PendingWriteNum int

	conns      map[net.Conn]struct{}
	mutexConns sync.Mutex
	wgConns    sync.WaitGroup

	initHandler ChannelInitializer
}

func (selft *Bootstrap) Handler(handler func(channel *Channel)) *Bootstrap {
	selft.initHandler = handler
	selft.conns = map[net.Conn]struct{}{}
	return selft
}

func (self *Bootstrap) RunServer() {
	l, err := net.Listen("tcp", self.Addr)
	defer l.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}

	var tempDelay time.Duration
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Printf("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		self.mutexConns.Lock()

		if len(self.conns) >= self.MaxConnNum {
			_ = conn.Close()
			self.mutexConns.Unlock()
			continue
		}

		self.wgConns.Add(1)

		self.conns[conn] = struct{}{}
		self.mutexConns.Unlock()

		go func(ch *Channel) {
			ch.runEventLoop()
			_ = ch.Close()

			self.mutexConns.Lock()
			delete(self.conns, conn)
			self.mutexConns.Unlock()

			self.wgConns.Done()
		}(self.initChannel(conn))
	}
}

func (b *Bootstrap) initChannel(conn net.Conn) *Channel {
	p := NewPipeline()

	c := NewChannel(conn, p)
	p.ch = c
	b.initHandler(c)

	p.fireNextChannelActive()
	return c
}
