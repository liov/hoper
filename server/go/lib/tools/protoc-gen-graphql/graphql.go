package main

import (
	"github.com/liov/hoper/server/go/lib/tools/clawer/protoc-gen-graphql/plugin"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{ParamFunc: plugin.Builder.Flags.Set}.Run(plugin.Generate)
}
