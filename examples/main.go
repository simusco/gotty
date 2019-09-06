package main

import (
	"github.com/golang/protobuf/proto"
	"simusco.com/gotty"
	"simusco.com/gotty/examples/live"
	"simusco.com/gotty/examples/service"
)

func main() {
	b := gotty.Bootstrap{
		Addr:            "127.0.0.1:9001",
		MaxConnNum:      100,
		PendingWriteNum: 100,
	}

	serviceHandler := &gotty.ServiceHandler{
		Services: map[int32]gotty.Service{
			1: &service.TestService{},
			2: &service.TestService{},
		},
	}

	protobufEncoder := &gotty.ProtobufEncoder{}
	protobufDecoder := &gotty.ProtobufDecoder{NewParam: func() proto.Message {
		return new(live.ParamVO)
	}}

	b.Handler(func(channel *gotty.Channel) {
		channel.Pipeline().
			AddLast(protobufEncoder).
			AddLast(protobufDecoder).
			AddLast(serviceHandler)
	}).RunServer()
}
