package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/mock/user.mock.go -package mock -source ../protobuf/user/user.service.pb.go UserServiceServer

func main() { run() }

var files = map[string][]string{
	"/utils/empty/*.gen.proto":      {"gogo_out=plugins=grpc"},
	"/utils/errorcode/errrep.proto": {"gogo_out=plugins=grpc"},
	"/utils/errorcode/*enum.proto":  {"enum_out=plugins=grpc"},
	"/utils/actor/message/*.proto":  {"gogo_out=plugins=grpc"},
	"/utils/response/*.gen.proto":   {"gogo_out=plugins=grpc"},
	"/utils/proto/gogo/*.gen.proto": {"gogo_out=plugins=grpc"},
	"/utils/proto/go/*.gen.proto":   {"go_out=plugins=grpc"},
	"/user/*service.proto": {"gogo_out=plugins=grpc",
		"grpc-gateway_out=logtostderr=true",
		"swagger_out=logtostderr=true",
		"govalidators_out=gogoimport=true",
		//"gogqlgen_out=gogoimport=false,paths=source_relative",
		"gqlgencfg_out=paths=source_relative",
		"gql_out=paths=source_relative"},
	"/user/*model.proto": {"gogo_out=plugins=grpc"},
	"/user/*enum.proto":  {"enum_out=plugins=grpc"},
	"/note/*service.proto": {"gogo_out=plugins=grpc",
		"grpc-gateway_out=logtostderr=true",
		"swagger_out=logtostderr=true",
		"govalidators_out=gogoimport=true"},
	"/note/*model.proto": {"gogo_out=plugins=grpc"},
}

var proto = flag.String("proto", "../../../../proto", "proto路径")
var env = []string{
	"GOARCH=" + runtime.GOARCH,
	"GOOS=" + runtime.GOOS,
}

func run() {
	pwd, _ := os.Getwd()
	*proto = pwd + "/" + *proto
	path := os.Getenv("GOPATH")
	include := "-I" + *proto + " -I" + path + "/src -I" + path + "/src/github.com/grpc-ecosystem/grpc-gateway -I" + path + "/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I" + path + "/src/github.com/gogo/protobuf/protobuf "

	for k, v := range files {
		for _, plugin := range v {
			arg := "protoc " + include + *proto + k + " --" + plugin + ",Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:" + pwd + "/protobuf"
			if strings.HasPrefix(plugin, "swagger_out") {
				arg = arg + "/api"
			}
			if strings.HasPrefix(plugin, "gql_out") || strings.HasPrefix(plugin, "gqlgencfg_out") {
				arg = arg + "/gql"
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
