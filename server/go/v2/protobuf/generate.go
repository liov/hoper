package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/liov/hoper/go/v2/utils/os"
	execi "github.com/liov/hoper/go/v2/utils/os/exec"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

func main() {
	//single("/content/moment.model.proto")
	run(*proto)
	//genutils(*proto+"/utils")
	//gengql()
}

const goOut = "go-patch_out=plugin=go,paths=source_relative"
const grpcOut = "go-patch_out=plugin=go-grpc,paths=source_relative"
const enumOut = "enum_out=plugins=grpc,paths=source_relative"
const gatewayOut = "grpc-gateway_out=logtostderr=true,paths=source_relative"
const openapiv2Out = "openapiv2_out=logtostderr=true"
const govalidatorsOut = "govalidators_out=gogoimport=true,paths=source_relative"
const gogoprotoOut = "gogo_out=plugins=grpc"
const gqlNogogoOut = "gqlgen_out=gogoimport=false,paths=source_relative"
const gqlOut = "graphql_out=paths=source_relative"

var service = []string{goOut, grpcOut,
	gatewayOut, openapiv2Out, govalidatorsOut,
	//gqlNogogoOut, gqlOut,
	//"gqlgencfg_out=paths=source_relative",
}

var model = []string{goOut, grpcOut}
var enum = []string{enumOut, goOut, grpcOut}

var gqlgen []string
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

var (
	proto                                         *string
	pwd, goList, gateway, protobuf, path, include string
)

func init() {
	proto = flag.String("proto", "../../../proto", "proto路径")
	pwd, _ = os.Getwd()
	*proto = pwd + "/" + *proto
	goList = `go list -m -f {{.Dir}} `
	gateway, _ = osi.CMD(
		goList + "github.com/grpc-ecosystem/grpc-gateway/v2",
	)

	//protopatch, _ := os2.CMD(goList + "github.com/liov/protopatch2")
	protobuf, _ = osi.CMD(goList + "google.golang.org/protobuf")
	//gogoProtoOut, _ := cmd.CMD(goList + "github.com/gogo/protobuf")
	path = os.Getenv("GOPATH")
	include = "-I" + *proto + " -I" + gateway + " -I" + gateway + "/third_party/googleapis -I" + protobuf + " -I" + path + "/src"
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

func single(file string) {
	for _, plugin := range model {
		arg := "protoc " + include + " " + *proto + file + " --" + plugin + ":" + pwd + "/protobuf"
		execi.Run(arg)
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
			arg := "protoc " + include + " " + dir + "/" + fileInfos[i].Name() + " --" + enumOut + ":" + pwd + "/protobuf"
			execi.Run(arg)
		}
		if strings.HasPrefix(dir, "/utils/proto/gogo/") {
			arg := "protoc -I" + *proto + " " + dir + " --gogo_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:" + pwd + "/protobuf"
			execi.Run(arg)
			continue
		}
		for _, plugin := range model {
			arg := "protoc " + include + " " + dir + "/*.proto" + " --" + plugin + ":" + pwd + "/protobuf"
			execi.Run(arg)
		}
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

const ymlTpl = `schema:
  - ./*.graphqls

# Where should the generated server code go?
exec:
  filename: ../../{{.}}/generated.gql.go
  package: {{.}}

# Enable Apollo federation support
federation:
  filename: ../../{{.}}/federation.gql.go
  package: {{.}}

model:
  filename: ../../{{.}}/models.gql.go
  package: {{.}}

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
