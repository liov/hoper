package main

import (
	"log"
	"os"
)

func main() {
	dir := "F:/Pictures"
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
		newDir := dir + "\\" + date
		_, err := os.Stat(newDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(newDir, 0666)
			if err != nil {
				log.Println(err)
			}
		}
		err = os.Rename(dir+"\\"+entity.Name(), newDir+"\\"+entity.Name())
		if err != nil {
			log.Println(err)
		}
	}
}
