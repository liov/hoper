package main

import (
	"golang.org/x/image/webp"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".webp") {
			file, err := os.Open(dir + "/" + f.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			image, err := webp.Decode(file)
			if err != nil {
				log.Println(err)
				continue
			}
			file.Close()
			jpegfile, err := os.Create(dir + "/" + f.Name() + ".jpg")
			if err != nil {
				log.Println(err)
				continue
			}
			err = jpeg.Encode(jpegfile, image, &jpeg.Options{Quality: 90})
			if err != nil {
				log.Println(err)
				continue
			}
			jpegfile.Close()
		}
	}
}
