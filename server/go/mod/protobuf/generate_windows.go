//go:build windows

package main

import (
	execi "github.com/actliboy/hoper/server/go/lib/utils/os/exec"
	"io/ioutil"
	"log"
	"strings"
)

var files = map[string][]string{
	/*
		"/utils/errorcode/errrep.proto": model,
		"/utils/errorcode/*enum.proto":  enum,
		"/utils/response/*.proto":       model,
		"/utils/request/*.proto":       model,
		"/utils/oauth/*.proto":          model,
		"/utils/time/*.proto":          model,
		"/utils/proto/gogo/*.gen.proto": {gogoprotoOut},
		"/utils/proto/go/*.proto":       {goOut},*/
	"/*service.proto": service,
	"/*model.proto":   model,
	"/*enum.proto":    enum,
}

func run(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if fileInfos[i].Name() == "utils" {
				continue
			}
			for k, v := range files {
				k = dir + "/" + fileInfos[i].Name() + k
				for _, plugin := range v {
					arg := "protoc " + include + " " + k + " --" + plugin + ":" + pwd + "/protobuf"
					if strings.HasPrefix(plugin, "openapiv2_out") {
						arg = arg + "/api"
					}
					if strings.HasPrefix(plugin, "graphql_out") || strings.HasPrefix(plugin, "gqlcfg_out") {
						arg = arg + "/gql"
					}
					//protoc-gen-gqlgen应该在最后生成，gqlgen会调用go编译器，protoc-gen-gqlgen会生成不存在的接口，编译不过去
					if strings.HasPrefix(plugin, "gqlgen_out") {
						gqlgen = append(gqlgen, arg)
						continue
					}
					execi.Run(arg)
				}
			}
			run(dir + "/" + fileInfos[i].Name())
		}
	}
}

func genutils(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			genutils(dir + "/" + fileInfos[i].Name())
		}
		if strings.HasSuffix(fileInfos[i].Name(), "enum.proto") {
			arg := "protoc " + include + " " + dir + "/*enum.proto" + " --" + enumOut + ":" + pwd + "/protobuf"
			execi.Run(arg)
			break
		}
	}

	for _, plugin := range model {
		arg := "protoc " + include + " " + dir + "/*.proto" + " --" + plugin + ":" + pwd + "/protobuf"
		execi.Run(arg)
	}
}
