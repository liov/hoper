package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/utils/os"
	execi "github.com/actliboy/hoper/server/go/lib/utils/os/exec"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

func main() {
	//single("/content/moment.model.proto")
	run(proto)
	genutils(proto + "/utils")
	//gengql()

}

const (
	goOut           = "go-patch_out=plugin=go,paths=source_relative"
	grpcOut         = "go-patch_out=plugin=go-grpc,paths=source_relative"
	enumOut         = "enum_out=plugins=grpc,paths=source_relative"
	gatewayOut      = "grpc-gin_out=paths=source_relative"
	openapiv2Out    = "openapiv2_out=logtostderr=true"
	govalidatorsOut = "govalidators_out=gogoimport=true,paths=source_relative"
	gogoprotoOut    = "gogo_out=plugins=grpc"
	gqlNogogoOut    = "gqlgen_out=gogoimport=false,paths=source_relative"
	gqlOut          = "graphql_out=paths=source_relative"
)

const (
	goListDir           = `go list -m -f {{.Dir}} `
	goListDep           = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis       = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	DepHoper            = "github.com/actliboy/hoper/server/go/lib"
	DepHoperWithVersion = "github.com/actliboy/hoper/server/go/lib@v0.0.0-20220713022058-3fd4e65cb7bc"
)

var (
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
	DepProtobuf    = "google.golang.org/protobuf"
)

var service = []string{goOut, grpcOut,
	gatewayOut, openapiv2Out, govalidatorsOut,
	//gqlNogogoOut, gqlOut,
	//"gqlgencfg_out=paths=source_relative",
}

var model = []string{goOut, grpcOut}
var enum = []string{enumOut, goOut}

var gqlgen []string

var (
	proto, genpath           string
	gopath, modPath, include string
)

func init() {
	pwd, _ := os.Getwd()
	protodef, _ := filepath.Abs("../../../proto")
	flag.StringVar(&proto, "proto", protodef, "proto路径")
	flag.StringVar(&genpath, "genpath", pwd+"/protobuf", "生成路径")
	stdPatch := flag.Bool("patch", false, "是否使用原生protopatch")
	flag.Parse()
	proto, _ = filepath.Abs(proto)
	genpath, _ = filepath.Abs(genpath)
	log.Println("proto:", proto)
	log.Println("genpath:", genpath)
	gopath = os.Getenv("GOPATH")
	if gopath != "" && !strings.HasSuffix(gopath, "/") {
		gopath = gopath + "/"
	}
	modPath = gopath + "pkg/mod/"

	generatePath := "generate" + time.Now().Format("150405")
	err := os.Mkdir(generatePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	generatePath = pwd + "/" + generatePath
	defer os.RemoveAll(generatePath)
	err = os.Chdir(generatePath)
	if err != nil {
		log.Fatal(err)
	}
	osi.CMD("go mod init generate")

	libHoperDir := getDepDir(DepHoper)
	if libHoperDir == "" {
		return
	} else {
		os.Chdir(libHoperDir)
		DepGrpcGateway, _ = osi.CMD(goListDep + DepGrpcGateway)
		DepProtopatch, _ = osi.CMD(goListDep + DepProtopatch)
		DepProtobuf, _ = osi.CMD(goListDep + DepProtobuf)
		os.Chdir(generatePath)
	}
	libGoogleDir := getDepDir(DepGoogleapis)

	libGatewayDir := getDepDir(DepGrpcGateway)

	protopatch := libHoperDir + "/protobuf"
	if *stdPatch {
		protopatch = getDepDir(DepProtopatch)
	}

	libProtobufDir := getDepDir(DepProtobuf)

	os.Chdir(pwd)
	//gogoProtoOut, _ := cmd.CMD(goListDir + "github.com/gogo/protobuf")
	include = "-I" + libGatewayDir + " -I" + protopatch +
		" -I" + libGoogleDir + " -I" + libHoperDir + "/protobuf -I" +
		libProtobufDir + " -I" + libHoperDir + "/protobuf/third" + " -I" + proto
	log.Println("include:", include)
}

func getDepDir(dep string) string {
	if !strings.Contains(dep, "@") {
		return modDepDir(dep)
	}
	depPath := modPath + dep
	_, err := os.Stat(depPath)
	if os.IsNotExist(err) {
		depPath = modDepDir(dep)
	}
	return depPath
}

func modDepDir(dep string) string {
	depPath, _ := osi.CMD(goListDir + dep)
	if depPath == "" {
		osi.CMD("go get " + dep)
		depPath, _ = osi.CMD(goListDir + dep)
	}
	return depPath
}

func single(path string) {
	for _, plugin := range model {
		arg := "protoc " + include + " " + path + " --" + plugin + ":" + genpath
		execi.Run(arg)
	}
}

func gengql() {
	gqldir := genpath + "/gql"
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

	for i := range gqlgen {
		words := osi.Split(gqlgen[i])
		cmd := exec.Command(words[0], words[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}
