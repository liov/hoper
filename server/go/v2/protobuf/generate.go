package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/liov/hoper/go/v2/utils/os2"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

func main() { run() }

const goOut = "go-patch_out=plugin=go,paths=source_relative"
const grpcOut = "go-patch_out=plugin=go-grpc,paths=source_relative"
const enumOut = "enum_out=plugins=grpc,paths=source_relative"
const gatewayOut = "grpc-gateway_out=logtostderr=true,paths=source_relative"
const openapiv2Out = "openapiv2_out=logtostderr=true"
const govalidatorsOut = "govalidators_out=gogoimport=true,paths=source_relative"
const gogoprotoOut = "gogo_out=plugins=grpc"

var gengql = true
var files = map[string][]string{
	"/utils/empty/*.proto":          {goOut, grpcOut},
	"/utils/errorcode/errrep.proto": {goOut, grpcOut},
	"/utils/errorcode/*enum.proto":  {enumOut, goOut, grpcOut},
	"/utils/actor/message/*.proto":  {goOut, grpcOut},
	"/utils/response/*.proto":       {goOut, grpcOut},
	"/utils/oauth/*.proto":          {goOut, grpcOut},
	"/utils/proto/gogo/*.gen.proto": {gogoprotoOut},
	"/utils/proto/go/*.proto":       {goOut},
	"/user/*service.proto": {goOut, grpcOut,
		gatewayOut,
		openapiv2Out,
		govalidatorsOut,
		"gqlgen_out=gogoimport=false,paths=source_relative",
		//"gqlgencfg_out=paths=source_relative",
		"graphql_out=paths=source_relative"},
	"/user/*model.proto": {goOut, grpcOut},
	"/user/*enum.proto":  {enumOut, goOut, grpcOut},
	"/note/*service.proto": {goOut, grpcOut,
		gatewayOut,
		openapiv2Out,
		govalidatorsOut},
	"/note/*model.proto": {goOut, grpcOut},
}

var proto = flag.String("proto", "../../../proto", "proto路径")

func run() {
	pwd, _ := os.Getwd()
	*proto = pwd + "/" + *proto
	goList := `go list -m -f {{.Dir}} `
	gateway, _ := os2.CMD(goList + "github.com/grpc-ecosystem/grpc-gateway/v2")
	protopatch, _ := os2.CMD(goList + "github.com/liov/protopatch2")
	protobuf, _ := os2.CMD(goList + "google.golang.org/protobuf")
	//gogoProtoOut, _ := cmd.CMD(goList + "github.com/gogo/protobuf")
	path := os.Getenv("GOPATH")
	include := "-I" + *proto + " -I" + gateway + " -I" + gateway + "/third_party/googleapis -I" + protopatch + " -I" + protobuf + " -I" + path + "/src"

	var gqlgen []string
	for k, v := range files {
		for _, plugin := range v {
			arg := "protoc " + include + " " + *proto + k + " --" + plugin + ":" + pwd + "/protobuf"
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
			if strings.HasPrefix(k, "/utils/proto/gogo/") {
				arg = "protoc -I" + *proto + " " + *proto + k + " --gogo_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:" + pwd + "/protobuf"
			}

			words := os2.Split(arg)
			cmd := exec.Command(words[0], words[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}
	if gengql {
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
				words := os2.Split(`gqlgen --verbose --config ` + config)
				cmd := exec.Command(words[0], words[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}
	}
	os.Chdir(pwd)

	for i := range gqlgen {
		words := os2.Split(gqlgen[i])
		cmd := exec.Command(words[0], words[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}

const ymlTpl = `schema:
  - ./*.graphqls

# Where should the generated server code go?
exec:
  filename: ../../{{.}}/generated.gql.go
  package: user

# Enable Apollo federation support
federation:
  filename: ../../{{.}}/federation.gql.go
  package: user

model:
  filename: ../../{{.}}/models.gql.go
  package: user

autobind:
  - "github.com/liov/hoper/go/v2/protobuf/{{.}}"
  - "github.com/liov/hoper/go/v2/protobuf/utils/response"
  - "github.com/liov/hoper/go/v2/protobuf/utils/oauth"

models:
  ID:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.UInt64
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
  Int32:
    model:
      - github.com/99designs/gqlgen/graphql.Int32
  Int64:
    model:
      - github.com/99designs/gqlgen/graphql.Int64
  Uint8:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Uint8
  Uint:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Uint
  Uint32:
      model:
        - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Uint32
  Uint64:
      model:
        - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Uint64
  Float32:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Float32
  Float64:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Float64
  Float:
    model:
      - github.com/99designs/gqlgen/graphql.Float
  Bytes:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.Bytes
  HttpResponse_HeaderEntry:
    model:
      - github.com/liov/hoper/go/v2/utils/net/http/api/graphql.HttpResponse_HeaderEntry
`

//经过一番查找，发现yaml语法对格式是非常严格的，不可以有制表符！不可以有制表符！不可以有制表符！
//缩进也有要求
