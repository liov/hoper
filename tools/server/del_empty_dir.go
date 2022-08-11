package main

import (
	"log"
	"os"
)

var CommonDir = "./"

const Sep = string(os.PathSeparator)

func main() {
	del(CommonDir)
}

func del(dir string) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	if len(fileInfos) == 0 {
		os.Remove(dir)
		log.Println("remove:", dir)
		return
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			del(dir + Sep + fileInfos[i].Name())
		}
	}
}
