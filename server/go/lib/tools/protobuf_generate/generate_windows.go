//go:build windows

package main

import (
	execi "github.com/liov/hoper/server/go/lib/utils/os/exec"
	"log"
	"os"
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
		"/utils/proto/go/*.proto":       {goOut},*/
	"/*service.proto": service,
	"/*.proto":        model,
	"/*enum.proto":    enum,
}

func run(dir string) {
	fileInfos, err := os.ReadDir(dir)
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
				protoc(v, k)
			}
			run(dir + "/" + fileInfos[i].Name())
		}
	}
}

func genutils(dir string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			genutils(dir + "/" + fileInfos[i].Name())
		}
		if strings.HasSuffix(fileInfos[i].Name(), "enum.proto") {
			arg := "protoc " + include + " " + dir + "/*enum.proto" + " --" + enumOut + ":" + genpath
			execi.Run(arg)
			break
		}
	}

	for _, plugin := range model {
		arg := "protoc " + include + " " + dir + "/*.proto" + " --" + plugin + ":" + genpath
		execi.Run(arg)
	}
}
