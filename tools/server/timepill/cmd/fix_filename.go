package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"os"
	"strings"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, nil)()
	special()
}

func normal() {
	dir := timepill.Conf.TimePill.PhotoPath + "/"
	dirInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range dirInfos {
		if dirInfos[i].IsDir() && strings.Contains(dirInfos[i].Name(), "-") {
			subDir := dir + dirInfos[i].Name() + "/"
			subDirs, err := os.ReadDir(subDir)
			if err != nil {
				log.Error(err)
			}
			for d := range subDirs {
				subDir := subDir + subDirs[d].Name() + "/"
				fileInfos, err := os.ReadDir(subDir)
				if err != nil {
					log.Error(err)
				}
				for j := range fileInfos {
					if !fileInfos[j].IsDir() && fileInfos[j].Name()[10] == '-' {
						log.Info(fileInfos[j].Name())
						path := subDir + fileInfos[j].Name()
						filename := fileInfos[j].Name()
						err = os.Rename(path, subDir+filename[0:10]+"_"+filename[11:])
						if err != nil {
							log.Error(err)
						}
					}
				}
			}

		}
	}
}

func special() {
	dir := timepill.Conf.TimePill.PhotoPath + "/0/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range files {
		subs := strings.Split(files[i].Name(), "_")
		err = os.Rename(dir+files[i].Name(), dir+subs[2]+subs[3])
		if err != nil {
			log.Error(err)
		}

	}

}
