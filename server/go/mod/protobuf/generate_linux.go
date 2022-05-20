//go:build linux

package main

import (
	execi "github.com/actliboy/hoper/server/go/lib/utils/os/exec"
	"io/ioutil"
	"log"
	"strings"
)

var stuffs = map[string][]string{
	/*
		"/utils/errorcode/errrep.proto": model,
		"/utils/errorcode/*enum.proto":  enum,
		"/utils/response/*.proto":       model,
		"/utils/request/*.proto":       model,
		"/utils/oauth/*.proto":          model,
		"/utils/time/*.proto":          model,
		"/utils/proto/gogo/*.gen.proto": {gogoprotoOut},
		"/utils/proto/go/*.proto":       {goOut},*/
	"service.proto": service,
	"model.proto":   model,
	"enum.proto":    enum,
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
			run(dir + "/" + fileInfos[i].Name())
			continue
		}
		for k, v := range stuffs {
			filename := fileInfos[i].Name()
			file := dir + "/" + filename
			if strings.HasSuffix(filename, k) {
				for _, plugin := range v {
					arg := "protoc " + include + " " + file + " --" + plugin + ":" + pwd + "/protobuf"
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
		}
	}
}

func genutils(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	for i := range fileInfos {
		filename := fileInfos[i].Name()
		if fileInfos[i].IsDir() {
			genutils(dir + "/" + filename)
			continue
		}
		if strings.HasSuffix(filename, "enum.proto") {
			arg := "protoc " + include + " " + dir + "/" + filename + " --" + enumOut + ":" + pwd + "/protobuf"
			execi.Run(arg)
			continue
		}
		single(dir + "/" + filename)
	}

}
