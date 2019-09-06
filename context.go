package gotty

import "log"

type HandlerContext struct {
	p       *Pipeline
	next    *HandlerContext
	prev    *HandlerContext
	handler Handler
}

func NewHandlerContext(p *Pipeline, handler Handler) *HandlerContext {
	return &HandlerContext{p: p, handler: handler}
}

func (h *HandlerContext) Write(data interface{}) {
	hc := h.findNextOutbound()
	if hc != nil {
		_ = hc.handler.(OutboundHandler).Write(hc, data)
	}
}

func (h *HandlerContext) Close() {
	hc := h.findNextOutbound()
	if hc != nil {
		_ = hc.handler.(OutboundHandler).Close(hc)
	}
}

func (h *HandlerContext) isInbound() bool {
	_, ok := h.handler.(InboundHandler)
	return ok
}

func (h *HandlerContext) isOutbound() bool {
	_, ok := h.handler.(OutboundHandler)
	return ok
}

func (h *HandlerContext) findNextInbound() *HandlerContext {
	next := h
	for {
		next = next.next
		if next.isInbound() {
			return next
		}
	}
}

func (h *HandlerContext) findNextOutbound() *HandlerContext {
	prev := h
	for {
		prev = prev.prev
		if prev.isOutbound() {
			return prev
		}
	}
}

func (h *HandlerContext) FireChannelActive() {
	hc := h.findNextInbound()
	if hc != nil {
		bb := hc.handler.(InboundHandler)
		err := bb.ChannelActive(hc)
		log.Printf("%v", err)
	}
}

func (h *HandlerContext) FireChannelRead(data interface{}) {
	hc := h.findNextInbound()
	if hc != nil {
		bb := hc.handler.(InboundHandler)
		err := bb.ChannelRead(hc, data)
		log.Printf("read error : %v", err)
	}
}
