package entity

import "google.golang.org/grpc"

type Character int

func (c *Character) Add(payload string, reply *string) error {
	test := "gg"
	reply = &test
	return nil
}

var Greeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Character.add",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/helloworld/helloworld/helloworld.proto",
}
