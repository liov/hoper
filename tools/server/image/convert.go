package main

import (
	"golang.org/x/image/webp"
	"image/jpeg"
	"log"
	"os"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	files, err := os.ReadDir(dir)
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
			os.Remove(dir + "/" + f.Name())
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
