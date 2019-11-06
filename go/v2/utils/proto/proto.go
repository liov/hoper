package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
)

var out  string

func main() {
	in:= flag.String("in","../../../protobuf","go protobuf")
	flag.StringVar(&out,"out","../../../../protobuf","通用 protobuf")
	flag.Parse()

}

func parse(in string)  {
	fileInfos,err := ioutil.ReadDir(in)
	if err!= nil {
		log.Error(err)
	}
	for i:= range fileInfos{
		fileName := in + string(os.PathSeparator) + fileInfos[i].Name()
		log.Info(fileName)
		if fileInfos[i].IsDir(){
			err := os.Mkdir(strings.Replace(fileName,in,out,1),0777)
			if err !=nil {
				log.Error(err)
			}
			parse(fileName)
		}else {
			replace(fileName)
		}
	}
}

func replace(in string)  {
	if path.Ext(in) == ".proto" {
		file1,err := os.Open(in)
		if err !=nil {
			log.Error(err)
		}
		file2,err := os.Create(strings.Replace(in,in,out,1))
		if err !=nil {
			log.Error(err)
		}
	}
}