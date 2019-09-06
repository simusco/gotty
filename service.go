package gotty

type Service interface {
	Execute(channel *Channel, data interface{})
}
