package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/os"
	execi "github.com/actliboy/hoper/server/go/lib/utils/os/exec"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

func main() {
	//single("/content/moment.model.proto")
	run(*proto)
	genutils(*proto + "/utils")
	//gengql()
	file, _ := os.Create(pwd + "/protobuf/build")
	defer file.Close()
	file.WriteString(time.Now().Format("2006-01-02 15:04:05"))
}

const goOut = "go-patch_out=plugin=go,paths=source_relative"
const grpcOut = "go-patch_out=plugin=go-grpc,paths=source_relative"
const enumOut = "enum_out=plugins=grpc,paths=source_relative"
const gatewayOut = "grpc-gin_out=paths=source_relative"
const openapiv2Out = "openapiv2_out=logtostderr=true"
const govalidatorsOut = "govalidators_out=gogoimport=true,paths=source_relative"
const gogoprotoOut = "gogo_out=plugins=grpc"
const gqlNogogoOut = "gqlgen_out=gogoimport=false,paths=source_relative"
const gqlOut = "graphql_out=paths=source_relative"

const googleapis = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"

var service = []string{goOut, grpcOut,
	gatewayOut, openapiv2Out, govalidatorsOut,
	//gqlNogogoOut, gqlOut,
	//"gqlgencfg_out=paths=source_relative",
}

var model = []string{goOut, grpcOut}
var enum = []string{enumOut, goOut}

var gqlgen []string

var (
	proto                                           *string
	pwd, goList, gateway, protobuf, gopath, include string
)

func init() {
	pwd, _ = os.Getwd()
	protodef, _ := filepath.Abs("../../../proto")
	proto = flag.String("proto", protodef, "proto路径")
	stdPatch := flag.Bool("patch", false, "是否使用原生protopatch")
	goList = `go list -m -f {{.Dir}} `

	google := gopath + "pkg/mod/" + googleapis
	_, err := os.Stat(google)
	if os.IsNotExist(err) {
		osi.CMD("go get " + googleapis)
		google, _ = osi.CMD(
			goList + "github.com/googleapis/googleapis",
		)
		osi.CMD("go mod tidy")
	}

	libDir, _ := osi.CMD(goList + "github.com/actliboy/hoper/server/go/lib")
	os.Chdir(libDir)
	gateway, _ = osi.CMD(
		goList + "github.com/grpc-ecosystem/grpc-gateway/v2",
	)

	protopatch := libDir + "/protobuf"
	if *stdPatch {
		protopatch, _ = osi.CMD(goList + "github.com/alta/protopatch")
	}
	protobuf, _ = osi.CMD(goList + "google.golang.org/protobuf")
	//gogoProtoOut, _ := cmd.CMD(goList + "github.com/gogo/protobuf")
	include = "-I" + gateway + " -I" + protopatch +
		" -I" + google + " -I" + libDir + "/protobuf -I" +
		protobuf + " -I" + libDir + "/protobuf/third" + " -I" + *proto
	os.Chdir(pwd)
}

func single(path string) {
	for _, plugin := range model {
		arg := "protoc " + include + " " + path + " --" + plugin + ":" + pwd + "/protobuf"
		execi.Run(arg)
	}
}

func gengql() {
	gqldir := pwd + "/protobuf/gql"
	fileInfos, err := ioutil.ReadDir(gqldir)
	if err != nil {
		log.Panicln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			os.Chdir(gqldir + "/" + fileInfos[i].Name())
			//这里用模板生成yml
			t := template.Must(template.New("yml").Parse(ymlTpl))
			config := fileInfos[i].Name() + `.service.gqlgen.yml`
			_, err := os.Stat(config)
			var file *os.File
			file, err = os.Create(config)
			if err != nil {
				log.Panicln(err)
			}
			t.Execute(file, fileInfos[i].Name())
			file.Close()
			execi.Run(`gqlgen --verbose --config ` + config)
		}
	}

	os.Chdir(pwd)

	for i := range gqlgen {
		words := osi.Split(gqlgen[i])
		cmd := exec.Command(words[0], words[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}
