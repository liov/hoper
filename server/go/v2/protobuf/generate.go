package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/mock/user.mock.go -package mock -source ../protobuf/user/user.service.pb.go UserServiceServer

func main() { run() }

var gengql = true
var files = map[string][]string{
	"/utils/empty/*.gen.proto":      {"gogo_out=plugins=grpc"},
	"/utils/errorcode/errrep.proto": {"gogo_out=plugins=grpc"},
	"/utils/errorcode/*enum.proto":  {"enum_out=plugins=grpc"},
	"/utils/actor/message/*.proto":  {"gogo_out=plugins=grpc"},
	"/utils/response/*.gen.proto":   {"gogo_out=plugins=grpc"},
	"/utils/oauth/*.gen.proto":      {"gogo_out=plugins=grpc"},
	"/utils/proto/gogo/*.gen.proto": {"gogo_out=plugins=grpc"},
	"/utils/proto/go/*.gen.proto":   {"go_out=plugins=grpc"},
	"/user/*service.proto": {"gogo_out=plugins=grpc",
		"grpc-gateway_out=logtostderr=true",
		"swagger_out=logtostderr=true",
		"govalidators_out=gogoimport=true",
		"gqlgen_out=gogoimport=false,paths=source_relative",
		//"gqlgencfg_out=paths=source_relative",
		"graphql_out=paths=source_relative"},
	"/user/*model.proto":   {"gogo_out=plugins=grpc"},
	"/user/*enum.proto":    {"enum_out=plugins=grpc"},
	"/user/*errcode.proto": {"enum_out=plugins=grpc"},
	"/note/*service.proto": {"gogo_out=plugins=grpc",
		"grpc-gateway_out=logtostderr=true",
		"swagger_out=logtostderr=true",
		"govalidators_out=gogoimport=true"},
	"/note/*model.proto": {"gogo_out=plugins=grpc"},
}

var proto = flag.String("proto", "../../../proto", "proto路径")

var env = []string{
	"GOARCH=" + runtime.GOARCH,
	"GOOS=" + runtime.GOOS,
}

func run() {
	pwd, _ := os.Getwd()
	*proto = pwd + "/" + *proto
	path := os.Getenv("GOPATH")
	include := "-I" + *proto + " -I" + path + "/src -I" + path + "/src/github.com/grpc-ecosystem/grpc-gateway -I" + path + "/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I" + path + "/src/github.com/gogo/protobuf/protobuf "

	var gqlgen []string
	for k, v := range files {
		for _, plugin := range v {
			arg := "protoc " + include + *proto + k + " --" + plugin + ",Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:" + pwd + "/protobuf"
			if strings.HasPrefix(plugin, "swagger_out") {
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
			if strings.HasPrefix(k, "/utils/proto/go/") {
				arg = "protoc -I" + *proto + " " + *proto + k + " --go_out=plugins=grpc:" + pwd + "/protobuf"
			}
			if strings.HasPrefix(k, "/utils/proto/gogo/") {
				arg = "protoc -I" + *proto + " " + *proto + k + " --gogo_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:" + pwd + "/protobuf"
			}
			words := split(arg)
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
				words := split(`gqlgen --verbose --config ` + config)
				cmd := exec.Command(words[0], words[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}
	}
	os.Chdir(pwd)

	for i := range gqlgen {
		words := split(gqlgen[i])
		cmd := exec.Command(words[0], words[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}

func split(line string) []string {
	var words []string
Words:
	for {
		line = strings.TrimLeft(line, " \t")
		if len(line) == 0 {
			break
		}
		if line[0] == '"' {
			for i := 1; i < len(line); i++ {
				c := line[i] // Only looking for ASCII so this is OK.
				switch c {
				case '\\':
					if i+1 == len(line) {
						log.Panic("bad backslash")
					}
					i++ // Absorb next byte (If it's a multibyte we'll get an error in Unquote).
				case '"':
					word, err := strconv.Unquote(line[0 : i+1])
					if err != nil {
						log.Panic("bad quoted string")
					}
					words = append(words, word)
					line = line[i+1:]
					// Check the next character is space or end of line.
					if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
						log.Panic("expect space after quoted argument")
					}
					continue Words
				}
			}
			log.Panic("mismatched quoted string")
		}
		i := strings.IndexAny(line, " \t")
		if i < 0 {
			i = len(line)
		}
		words = append(words, line[0:i])
		line = line[i:]
	}
	// Substitute command if required.

	// Substitute environment variables.
	for i, word := range words {
		words[i] = os.Expand(word, expandVar)
	}
	return words
}

func expandVar(word string) string {
	w := word + "="
	for _, e := range env {
		if strings.HasPrefix(e, w) {
			return e[len(w):]
		}
	}
	return os.Getenv(word)
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
