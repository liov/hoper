package main

import (
	"fmt"
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"log"
	"os"
	"strings"
)

func main() {
	transferCodec("F:\\B站\\")
}

const command = `ffmpeg  -hwaccel qsv -c:v h264_qsv -i %s -c:v hevc_qsv -y -gpu_copy on -c:a copy %s.mp4`

const cmd1 = `preset=veryslow,profile=main,look_ahead=1,global_quality=18`

func transferCodec(commondir string) {

	files, _ := os.ReadDir(commondir)

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".flv") {
			continue
		}
		filePath := commondir + file.Name()
		command := fmt.Sprintf(command, filePath, filePath[:len(filePath)-len(".flv")])
		log.Println(command)
		res, err := osi.CMD(command)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(res)
		//log.Println(os.Remove(filePath))
	}

}
