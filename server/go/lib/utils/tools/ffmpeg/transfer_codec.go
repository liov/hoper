package ffmpeg

import (
	"fmt"
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"log"
)

const param = "-global_quality 20"

const TransferCodecCmd = `ffmpeg  -hwaccel_output_format qsv -c:v h264_qsv -i %s -c:v hevc_qsv -preset veryslow -g 60 -gpu_copy 1 -c:a copy %s`

const cmd1 = `preset=veryslow,profile=main,look_ahead=1,global_quality=18`

func H264ToH265(filePath, dst string) error {
	command := fmt.Sprintf(TransferCodecCmd, filePath, dst)
	log.Println(command)
	res, err := osi.CMD(command)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	return nil
}
