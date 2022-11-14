package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/tools/ffmpeg"
	"log"
	"os"
	"strings"
)

func main() {
	transferFormat("F:\\B站")
}

func transferCodec(commondir string) {
	transfer(commondir, ffmpeg.H264ToH265)
}

func transferFormat(commondir string) {
	transfer(commondir, ffmpeg.TransferFormat)
}

func transfer(commondir string, call func(string, string) error) {

	files, _ := os.ReadDir(commondir)

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".flv") {
			continue
		}
		filePath := commondir + fs.PathSeparator + file.Name()
		err := call(filePath, filePath[:len(filePath)-len("flv")]+"mp4")
		if err != nil {
			log.Println(err)
		}
	}

}
