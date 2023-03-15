package main

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/tools/protoc-gen-enum/plugin"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"os"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	/*	file, _ := os.Create("input")
		file.Write(input)
		file.Close()*/
	var request pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(input, &request)
	if err != nil {
		panic(err)
	}

	opts := protogen.Options{}

	builder, err := plugin.New(opts, &request)
	if err != nil {
		panic(err)
	}

	response, err := builder.Generate()
	if err != nil {
		panic(err)
	}

	out, err := proto.Marshal(response)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(out))
	/*	var (
		flags   flag.FlagSet
	)*/
	// protogen.Options{ParamFunc: flags.Set}.Run(plugin.Generate)
}
