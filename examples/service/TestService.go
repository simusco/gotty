package service

import "simusco.com/gotty"

type TestService struct {
}

func (t *TestService) Execute(context *gotty.HandlerContext, data interface{}) {

	context.Write(nil)

}
