package gotty

import (
	"errors"
	"log"
	"simusco.com/gotty/examples/live"
)

type InboundHandler interface {
	ChannelRead(c *HandlerContext, data interface{}) error
	ChannelActive(c *HandlerContext) error
	Handler
}

type OutboundHandler interface {
	Write(c *HandlerContext, data interface{}) error
	Close(c *HandlerContext) error
	Handler
}

type Handler interface {
	ErrorCaught(c *HandlerContext, err error)
}

type ServiceHandler struct {
	Services map[int32]Service
}

func (sh *ServiceHandler) ChannelRead(c *HandlerContext, data interface{}) error {
	vo := data.(*live.ParamVO)
	if service, ok := sh.Services[vo.Event]; ok {
		service.Execute(c.p.ch, data)
		return nil
	}
	return errors.New("找不到服务实现")
}

func (sh *ServiceHandler) ChannelActive(c *HandlerContext) error {
	log.Printf("tcp_handler call")
	return nil
}

func (sh *ServiceHandler) ErrorCaught(c *HandlerContext, err error) {
}
