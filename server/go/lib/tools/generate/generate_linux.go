//go:build linux

package main

import (
	execi "github.com/liov/hoper/server/go/lib/utils/os/exec"
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
	".proto":        model,
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
				protoc(v, file)
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
			arg := "protoc " + include + " " + dir + "/" + filename + " --" + enumOut + ":" + genpath
			execi.Run(arg)
			continue
		}
		single(dir + "/" + filename)
	}

}
