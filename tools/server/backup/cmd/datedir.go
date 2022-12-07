package main

import (
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"log"
	"os"
)

func main() {
	dir := "F:\\Pictures\\pron\\pic\\WeiXin"
	entities, _ := os.ReadDir(dir)
	if len(entities) == 0 {
		return
	}

	for _, entity := range entities {
		if entity.IsDir() {
			continue
		}
		info, _ := entity.Info()
		date := info.ModTime().Format("200601")
		newDir := dir + fs.PathSeparator + date
		_, err := os.Stat(newDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(newDir, 0666)
			if err != nil {
				log.Println(err)
			}
		}
		oldpath := dir + fs.PathSeparator + entity.Name()
		newPath := newDir + fs.PathSeparator + entity.Name()
		_, err = os.Stat(newPath)
		if os.IsNotExist(err) {
			log.Println("rename:", newPath)
			err = os.Rename(oldpath, newPath)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("delete:", oldpath)
			os.Remove(oldpath)
		}

	}
}
