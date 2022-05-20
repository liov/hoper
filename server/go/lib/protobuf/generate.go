package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/utils/os"
	execi "github.com/actliboy/hoper/server/go/lib/utils/os/exec"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

func main() {
	//single("/content/moment.model.proto")
	genutils(proto)
	//gengql()
	os.Chdir(pwd)
}

const goOut = "go-patch_out=plugin=go,paths=source_relative"
const grpcOut = "go-patch_out=plugin=go-grpc,paths=source_relative"
const enumOut = "enum_out=plugins=grpc,paths=source_relative"

const googleapis = "googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32"

var model = []string{goOut, grpcOut}

var (
	proto                                           string
	pwd, goList, gateway, protobuf, gopath, include string
)

func init() {
	gopath = os.Getenv("GOPATH")
	if strings.HasSuffix(gopath, "/") {
		gopath = gopath[:len(gopath)-1]
	}
	stdPatch := flag.Bool("patch", false, "是否使用原生protopatch")
	pwd, _ = os.Getwd()
	proto = pwd + "/protobuf"
	goList = `go list -m -f {{.Dir}} `
	gateway, _ = osi.CMD(
		goList + "github.com/grpc-ecosystem/grpc-gateway/v2",
	)
	google := gopath + "pkg/mod/github.com/" + googleapis
	_, err := os.Stat(google)
	if os.IsNotExist(err) {
		osi.CMD("go get " + googleapis)
		google, _ = osi.CMD(
			goList + "github.com/googleapis/googleapis",
		)
		osi.CMD("go mod tidy")
	}
	protopatch := proto
	if *stdPatch {
		protopatch, _ = osi.CMD(goList + "github.com/alta/protopatch")
	}
	protobuf, _ = osi.CMD(goList + "google.golang.org/protobuf")
	//gogoProtoOut, _ := cmd.CMD(goList + "github.com/gogo/protobuf")

	include = "-I" + gateway + " -I" + google + " -I" +
		protobuf + " -I" + protopatch + " -I" + gopath + "/src" + " -I" + proto
}

func genutils(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Contains(dir, "lib/protobuf/third") {
		return
	}

	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			genutils(dir + "/" + fileInfos[i].Name())
		}
		if strings.HasSuffix(fileInfos[i].Name(), "enum.proto") {
			arg := "protoc " + include + " " + dir + "/" + fileInfos[i].Name() + " --" + enumOut + ":" + proto
			execi.Run(arg)
		}
	}
	if strings.Contains(dir, "lib/protobuf/utils/gogo") {
		if strings.HasSuffix(dir, ".gen.proto") {
			arg := "protoc -I" + proto + " " + dir + "/*.gen.proto --gogo_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:" + proto
			execi.Run(arg)
		}
		return
	}

	for _, plugin := range model {
		arg := "protoc " + include + " " + dir + "/*.proto --" + plugin + ":" + proto
		execi.Run(arg)
	}
}
