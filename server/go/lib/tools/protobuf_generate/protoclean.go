package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var in, out string

func clean() {
	in, _ = filepath.Abs(proto)
	out += "/_std"
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
		data, err := ioutil.ReadFile(src)
		if err != nil {
			log.Println(err)
		}
		newFilePath := strings.Replace(src, in, out, 1)
		reg := regexp.MustCompile(`import "github.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import "protoc-gen-openapiv2.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import "utils/enum/enum.proto.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import "utils/graphql/graphql.proto.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile(`import "google/api/annotations.proto.*\n`)
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile("\\[\\([\\w.]*\\)[=\"/:\\\\;\n\\-('{},\u4e00-\u9fa5\uff0c]*]")
		data = reg.ReplaceAll(data, nil)
		reg = regexp.MustCompile("option \\([\\w\\s.)/\\[\\]=\":*@\n\\-('{},\u4e00-\u9fa5]*;")
		data = reg.ReplaceAll(data, nil)
		err = ioutil.WriteFile(newFilePath, data, 0666)
		if err != nil {
			log.Println(err)
		}
	}
}
