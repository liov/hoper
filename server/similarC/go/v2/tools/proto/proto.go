package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
)

var in,out string

func main() {
	flag.StringVar(&in,"in", "../protobuf", "go protobuf")
	flag.StringVar(&out, "out", "../../../proto", "通用 protobuf")
	flag.Parse()
	_, err := os.Stat(out)
	if os.IsNotExist(err){
		os.Mkdir(out,os.ModePerm)
	}
	parse(in)
}

func parse(src string) {
	fileInfos, err := ioutil.ReadDir(src)
	if err != nil {
		log.Error(err)
	}
	for i := range fileInfos {
		fileName := src + "/" + fileInfos[i].Name()
		if fileInfos[i].IsDir() {
			newPath :=strings.Replace(fileName, in, out, 1)
			_, err := os.Stat(newPath)
			if os.IsNotExist(err){
				err = os.Mkdir(newPath, os.ModePerm)
				if err != nil {
					log.Error(err)
				}
			}
			parse(fileName)
		} else {
			log.Info(fileName)
			replace(fileName)
		}
	}
}

func replace(src string) {
	if path.Ext(src) == ".proto" {
		data, err := ioutil.ReadFile(src)
		if err != nil {
			log.Error(err)
		}
		newFilePath := strings.Replace(src, in, out, 1)
		reg := regexp.MustCompile(`import \"github.com/gogo/protobuf/gogoproto.*\n`)
		data = reg.ReplaceAll([]byte(data), nil)
		reg = regexp.MustCompile(`\[.*\]`)
		data = reg.ReplaceAll([]byte(data), nil)
		reg = regexp.MustCompile(`option \(gogoproto.*\n`)
		data = reg.ReplaceAll(data, nil)
		err =ioutil.WriteFile(newFilePath,data,0666)
		if err != nil {
			log.Error(err)
		}
	}
}
