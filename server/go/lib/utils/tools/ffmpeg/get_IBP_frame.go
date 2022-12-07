package ffmpeg

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	osi "github.com/liov/hoper/server/go/lib/utils/os"
	"log"
)

const GetFrameCmd = `ffmpeg -i %s -vf "select=eq(pict_type\,%s)" -fps_mode vfr -qscale:v 2 -f image2 %s/%%03d.jpg`

type Frame int

func (f Frame) String() string {
	switch f {
	case I:
		return "I"
	case P:
		return "P"
	case B:
		return "B"
	}
	return "unknown"
}

const (
	I Frame = iota
	P
	B
)

func GetFrame(src string, f Frame) error {
	//cmd := `ffmpeg -i ` + src + ` -vf "select=eq(pict_type\,` + f.String() + `)" -vsync vfr -qscale:v 2 -f image2 ` + dst + `/%03d.jpg`
	dst := fs.GetDir(src) + f.String() + "Frame"
	fs.Mkdir(dst)
	cmd := fmt.Sprintf(GetFrameCmd, src, f.String(), dst)
	res, err := osi.QuotedCMD(cmd)
	log.Println(res)
	return err
}
