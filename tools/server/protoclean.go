package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

var in, out string

func main() {
	log.SetFlags(15)
	flag.StringVar(&in, "in", "../../../proto", "go protobuf")
	flag.StringVar(&out, "out", "../../../proto_std", "通用 protobuf")
	flag.Parse()
	_, err := os.Stat(out)
	if os.IsNotExist(err) {
		os.Mkdir(out, os.ModePerm)
	}
	parse(in)
}

func parse(src string) {
	fileInfos, err := ioutil.ReadDir(src)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		fileName := src + "/" + fileInfos[i].Name()
		if fileInfos[i].IsDir() {
			///utils/proto文件夹不需要过去
			if strings.HasSuffix(fileName, "utils/proto") {
				continue
			}
			newPath := strings.Replace(fileName, in, out, 1)
			_, err := os.Stat(newPath)
			if os.IsNotExist(err) {
				err = os.Mkdir(newPath, os.ModePerm)
				if err != nil {
					log.Println(err)
				}
			}
			parse(fileName)
		} else {
			log.Println(fileName)
			replace(fileName)
		}
	}
}

func replace(src string) {
	if path.Ext(src) == ".proto" {
		if strings.HasSuffix(src, ".pub.proto") || strings.HasSuffix(src, ".imp.proto") {
			return
		}
		data, err := ioutil.ReadFile(src)
		if err != nil {
			log.Println(err)
		}
		newFilePath := strings.Replace(src, in, out, 1)
		reg := regexp.MustCompile(`import \"github.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import \"protoc-gen-openapiv2.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import \"utils/proto/gogo/enum.proto.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import \"utils/proto/gogo/graphql.proto.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import \"google/api/annotations.proto.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile("\\[\\([\\w\\s\\.\\)\\=\"/\\:\\\\;\n\\-\\(\\'\\{\\}\\,\u4e00-\u9fa5\uff0c]*\\]")
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile("option \\([\\w\\s\\.\\)/\\[\\]=\":*@\n\\-\\('\\{\\}\\,\u4e00-\u9fa5]*;")
		data = reg.ReplaceAll(data, nil)
		data = bytes.ReplaceAll(data, []byte(".pub.proto"), []byte(".proto"))
		data = bytes.ReplaceAll(data, []byte(".imp.proto"), []byte(".gen.proto"))
		err = ioutil.WriteFile(newFilePath, data, 0666)
		if err != nil {
			log.Println(err)
		}
	}
}
