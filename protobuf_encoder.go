package gotty

type ProtobufEncoder struct {
}

func (ProtobufEncoder) Write(c *HandlerContext, data interface{}) error {
	panic("implement me")
}

func (ProtobufEncoder) ErrorCaught(c *HandlerContext, err error) {
	panic("implement me")
}
