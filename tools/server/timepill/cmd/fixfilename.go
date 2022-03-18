package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"os"
	"strings"
)

func main() {
	dir := "/data/timepill/"
	dirInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range dirInfos {
		if dirInfos[i].IsDir() {
			subDir := dir + dirInfos[i].Name() + "/"
			fileInfos, err := os.ReadDir(subDir)
			if err != nil {
				log.Error(err)
			}
			for j := range fileInfos {
				if !fileInfos[j].IsDir() && strings.Contains(fileInfos[j].Name(), " ") {
					log.Info(fileInfos[j].Name())
					path := subDir + fileInfos[j].Name()
					filename := fileInfos[j].Name()
					err = os.Rename(path, subDir+filename[0:10]+filename[11:])
					if err != nil {
						log.Error(err)
					}
				}
			}

		}
	}
}
