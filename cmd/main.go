package main

import (
	"github.com/EvgenyOvsov/protoc-gateway-impl/pkg/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	generator := generator.NewGenerator()

	proto := protogen.Options{
		ParamFunc: generator.SetParamsFunc,
	}
	proto.Run(generator.Handle)
}
