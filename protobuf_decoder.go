package gotty

import (
	"github.com/golang/protobuf/proto"
	"log"
)

type ProtobufDecoder struct {
	NewParam func() proto.Message
}

func (self *ProtobufDecoder) ChannelActive(c *HandlerContext) error {
	log.Println("channel active")
	return nil
}

func (self *ProtobufDecoder) ChannelRead(c *HandlerContext, data interface{}) error {
	if bytes, ok := data.([]byte); ok {

		params := self.NewParam()
		err := proto.Unmarshal(bytes, params)

		if err != nil {
			self.ErrorCaught(c, err)
		}

		c.FireChannelRead(params)
	}
	return nil
}

func (self *ProtobufDecoder) ErrorCaught(c *HandlerContext, err error) {
	panic("implement me")
}
