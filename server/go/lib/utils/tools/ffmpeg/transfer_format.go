package ffmpeg

import (
	"fmt"
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"log"
)

const TransferFormatCmd = `ffmpeg -hwaccel qsv -i %s -c copy -y %s`

func TransferFormat(filePath, dst string) error {
	command := fmt.Sprintf(TransferFormatCmd, filePath, dst)
	log.Println(command)
	res, err := osi.CMD(command)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	return nil
}
