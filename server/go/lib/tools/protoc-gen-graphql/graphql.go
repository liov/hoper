package main

import (
	"github.com/actliboy/hoper/server/go/lib/tools/protoc-gen-graphql/plugin"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{ParamFunc: plugin.Builder.Flags.Set}.Run(plugin.Generate)
}
