package image

import (
	"golang.org/x/image/webp"
	"image/jpeg"
	"os"
	"path"
)

func WebpToJpg(src string, quality int) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	image, err := webp.Decode(file)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	err = os.Remove(src)
	if err != nil {
		return err
	}
	jpegfile, err := os.Create(path.Dir(src) + ".jpg")
	if err != nil {
		return err
	}
	err = jpeg.Encode(jpegfile, image, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}
	return jpegfile.Close()
}
