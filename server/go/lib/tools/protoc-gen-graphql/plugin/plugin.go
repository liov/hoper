package plugin

import (
	"flag"
	"google.golang.org/protobuf/compiler/protogen"
)

type builder struct {
	plugin *protogen.Plugin
	Flags  flag.FlagSet
}

var Builder = &builder{}

func Generate(gen *protogen.Plugin) error {
	genFileMap := make(map[string]*protogen.GeneratedFile)

	for _, protoFile := range gen.Files {
		fileName := protoFile.GeneratedFilenamePrefix + ".pb.gql.go"
		g := gen.NewGeneratedFile(fileName, ".")
		genFileMap[fileName] = g

	}

	for _, protoFile := range gen.Files {
		fileName := protoFile.GeneratedFilenamePrefix + ".pb.gql.go"
		g, ok := genFileMap[fileName]
		if !ok || protoFile.GeneratedFilenamePrefix == "google/protobuf/descriptor" {
			g.Skip()
			continue
		}

		g.P("package ", protoFile.GoPackageName)

	}

	return nil
}
