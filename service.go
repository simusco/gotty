package gotty

type Service interface {
	Execute(c *HandlerContext, data interface{})
}
