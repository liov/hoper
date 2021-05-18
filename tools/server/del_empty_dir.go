package main

import (
	"log"
	"os"

	"tools/pro"
)

func main() {
	del(pro.CommonDir)
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
			del(dir + pro.Sep + fileInfos[i].Name())
		}
	}
}
