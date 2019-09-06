package gotty

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Channel struct {
	sync.Mutex
	conn     net.Conn
	pipeline *Pipeline
	isClose  bool
}

func NewChannel(conn net.Conn, pipeline *Pipeline) *Channel {
	ch := &Channel{}
	ch.isClose = false
	ch.pipeline = pipeline
	ch.conn = conn
	return ch
}

func (ch *Channel) Pipeline() *Pipeline {
	return ch.pipeline
}

func (ch *Channel) runEventLoop() {
	for {
		if ch.isClose {
			log.Printf("连接已经关闭")
			break
		}

		var bytes2 []byte
		_, err := ch.doRead(bytes2)
		if err != nil {
			_ = ch.Close()
			fmt.Println("客户端断开了")
			break
		}

		ch.pipeline.fireNextChannelRead(bytes2)
	}
}

func (ch *Channel) doRead(p []byte) (n int, err error) {
	if ch.isClose {
		return -1, errors.New("连接已经关闭")
	}

	buf := make([]byte, 4)
	if _, err := io.ReadFull(ch.conn, buf); err != nil {
		return -1, err
	}

	msgLen := binary.BigEndian.Uint32(buf)

	if msgLen > 4096 {
		return -1, errors.New("message to long")
	}

	msg := make([]byte, msgLen)
	if _, err := io.ReadFull(ch.conn, msg); err != nil {
		return -1, err
	}

	p = msg

	return -1, nil
}

func (ch *Channel) doWrite(p []byte) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(len(p)))

	buf := new(bytes.Buffer)
	buf.Write(b)
	buf.Write(p)
}

func (ch *Channel) Close() error {
	if ch.isClose {
		return nil
	}

	err := ch.conn.Close()

	ch.isClose = true

	return err
}
