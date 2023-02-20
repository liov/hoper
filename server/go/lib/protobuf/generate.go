package main

import (
	"flag"
	execi "github.com/liov/hoper/server/go/lib/utils/os/exec"
	_go "github.com/liov/hoper/server/go/lib/utils/tools/gocmd"
	"log"
	"os"
	"strings"
)

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

func main() {
	//single("/content/moment.model.proto")
	genutils(proto)
	//gengql()
	os.Chdir(pwd)
}

const goOut = "go-patch_out=plugin=go,paths=source_relative"
const grpcOut = "go-patch_out=plugin=go-grpc,paths=source_relative"
const enumOut = "enum_out=paths=source_relative"

const (
	goListDir     = `go list -m -f {{.Dir}} `
	goListDep     = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	DepHoper      = "github.com/liov/hoper/server/go/lib"
)

var (
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
)

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

	libGatewayDir := _go.GetDepDir(DepGrpcGateway)
	libGoogleDir := _go.GetDepDir(DepGoogleapis)
	protopatch := proto
	if *stdPatch {
		protopatch = _go.GetDepDir(DepProtopatch)
	}

	include = "-I" + libGatewayDir + " -I" + libGoogleDir + " -I" + protopatch + " -I" + gopath + "/src" + " -I" + proto
}

func genutils(dir string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Contains(dir, "protobuf/third") {
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

	for _, plugin := range model {
		arg := "protoc " + include + " " + dir + "/*.proto --" + plugin + ":" + proto
		execi.Run(arg)
	}

}
