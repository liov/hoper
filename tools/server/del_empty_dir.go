package main

import (
	"log"
	"os"
)

func main() {
	del(por.CommonDir)
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
			del(dir + por.Sep + fileInfos[i].Name())
		}
	}
}
