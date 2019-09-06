package gotty

type Pipeline struct {
	head *HandlerContext
	tail *HandlerContext
	ch   *Channel
}

func NewPipeline() *Pipeline {
	p := &Pipeline{}
	p.tail = &HandlerContext{p, nil, nil, &tailHandler{}}
	p.head = &HandlerContext{p, nil, nil, &headHandler{}}
	p.head.next = p.tail
	p.tail.prev = p.head
	return p
}

func (p *Pipeline) AddLast(handler Handler) *Pipeline {
	prev := p.tail.prev
	newH := NewHandlerContext(p, handler)
	newH.prev = prev
	newH.next = p.tail
	prev.next = newH
	p.tail.prev = newH
	return p
}

func (pipeline *Pipeline) fireNextChannelActive() {
	pipeline.head.FireChannelActive()
}

func (pipeline *Pipeline) fireNextChannelRead(data interface{}) {
	pipeline.head.FireChannelRead(data)
}

type headHandler struct {
}

func (h *headHandler) ChannelRead(c *HandlerContext, data interface{}) error {
	c.FireChannelRead(data)
	return nil
}

func (h *headHandler) ChannelActive(c *HandlerContext) error {
	c.FireChannelActive()
	return nil
}

func (h *headHandler) ErrorCaught(c *HandlerContext, err error) {

}

func (h *headHandler) Write(c *HandlerContext, data interface{}) error {
	b, ok := data.([]byte)
	if ok {
		c.p.ch.doWrite(b)
	}
	return nil
}

func (h *headHandler) Close(c *HandlerContext) error {
	return c.p.ch.Close()
}

type tailHandler struct {
}

func (t *tailHandler) ChannelRead(c *HandlerContext, data interface{}) error {
	return nil
}

func (t *tailHandler) ChannelActive(c *HandlerContext) error {
	return nil
}

func (t *tailHandler) ErrorCaught(c *HandlerContext, err error) {

}

func (t *tailHandler) Write(c *HandlerContext, data interface{}) error {
	c.Write(data)
	return nil
}

func (t *tailHandler) Close(c *HandlerContext) error {
	c.Close()
	return nil
}
