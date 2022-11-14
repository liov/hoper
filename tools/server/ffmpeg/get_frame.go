package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/tools/ffmpeg"
	"log"
)

func main() {
	err := ffmpeg.GetFrame(`F:\Pictures\pron\baoyu\202107\11758_1627484824.mp4`, ffmpeg.I)
	if err != nil {
		log.Println(err)
	}
}
