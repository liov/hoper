package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	fixfilename("D:/F/timepill/2022-01-27/")
}

func fixfilename(dir string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && strings.HasSuffix(fileInfos[i].Name(), "!large") {
			os.Rename(dir+fileInfos[i].Name(), dir+strings.TrimSuffix(fileInfos[i].Name(), "!large"))
		}
	}
}
