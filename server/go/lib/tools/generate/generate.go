package main

import (
	_go "github.com/actliboy/hoper/server/go/lib/utils/go"
	"github.com/actliboy/hoper/server/go/lib/utils/os"
	execi "github.com/actliboy/hoper/server/go/lib/utils/os/exec"
	"github.com/spf13/cobra"
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

var rootCmd = &cobra.Command{
	Use: "generate",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		getInclude()
	},
}

func main() {
	//single("/content/moment.model.proto")
	rootCmd.Execute()
	//gengql()

}

const (
	goOut           = "go-patch_out=plugin=go,paths=source_relative"
	grpcOut         = "go-patch_out=plugin=go-grpc,paths=source_relative"
	enumOut         = "enum_out=paths=source_relative"
	gatewayOut      = "grpc-gin_out=paths=source_relative"
	openapiv2Out    = "openapiv2_out=logtostderr=true"
	govalidatorsOut = "govalidators_out=paths=source_relative"
	gqlNogogoOut    = "gqlgen_out=paths=source_relative"
	gqlOut          = "graphql_out=paths=source_relative"
	dartOut         = "dart_out=grpc"
)

const (
	goListDir     = `go list -m -f {{.Dir}} `
	goListDep     = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	DepHoper      = "github.com/actliboy/hoper/server/go/lib"
)

var (
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
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
	proto, genpath string
	include        string
	stdPatch       bool
)

func init() {
	protodef, _ := filepath.Abs("../../../proto")
	pwd, _ := os.Getwd()
	pflag := rootCmd.PersistentFlags()
	pflag.StringVarP(&proto, "proto", "p", protodef, "proto file path")
	pflag.StringVarP(&genpath, "genpath", "g", pwd+"/protobuf", "generate path")
	pflag.BoolVar(&stdPatch, "patch", false, "是否使用原生protopatch")
	rootCmd.AddCommand(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "go",
		Run: func(cmd *cobra.Command, args []string) {
			run(proto)
			genutils(proto + "/utils")
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "dart",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "ts",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})

}

func getInclude() {
	pwd, _ := os.Getwd()
	proto, _ = filepath.Abs(proto)
	genpath, _ = filepath.Abs(genpath)
	log.Println("proto:", proto)
	log.Println("genpath:", genpath)
	_, err := os.Stat(genpath + "/api")
	if os.IsNotExist(err) {
		err = os.Mkdir(genpath+"/api", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	generatePath := "generate" + time.Now().Format("150405")
	err = os.Mkdir(generatePath, os.ModePerm)
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

	libHoperDir := _go.GetDepDir(DepHoper)

	if libHoperDir == "" {
		return
	} else {
		os.Chdir(libHoperDir)
		DepGrpcGateway, _ = osi.CMD(goListDep + DepGrpcGateway)
		DepProtopatch, _ = osi.CMD(goListDep + DepProtopatch)
		os.Chdir(generatePath)
	}
	libGoogleDir := _go.GetDepDir(DepGoogleapis)

	libGatewayDir := _go.GetDepDir(DepGrpcGateway)

	protopatch := libHoperDir + "/protobuf"
	if stdPatch {
		protopatch = _go.GetDepDir(DepProtopatch)
	}

	os.Chdir(pwd)

	include = "-I" + libGatewayDir + " -I" + protopatch +
		" -I" + libGoogleDir + " -I" + libHoperDir + "/protobuf -I" + libHoperDir + "/protobuf/third  -I" + proto
	log.Println("include:", include)

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

func protoc(plugins []string, file string) {
	for _, plugin := range plugins {
		arg := "protoc " + include + " " + file + " --" + plugin + ":" + genpath
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
