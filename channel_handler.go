package gotty

import (
	"errors"
	"log"
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
	Services     map[int32]Service
	GetEventCode func(data interface{}) int32
}

func (sh *ServiceHandler) ChannelRead(c *HandlerContext, data interface{}) error {
	event := sh.GetEventCode(data)

	if event == -1 {
		return errors.New("丢失了Event参数")
	}

	if service, ok := sh.Services[event]; ok {
		service.Execute(c, data)
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
