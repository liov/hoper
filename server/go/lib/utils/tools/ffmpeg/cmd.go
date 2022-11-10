package ffmpeg

import (
	"fmt"
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
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

func GetFrame(src, dst string, f Frame) error {
	//cmd := `ffmpeg -i ` + src + ` -vf "select=eq(pict_type\,` + f.String() + `)" -vsync vfr -qscale:v 2 -f image2 ` + dst + `/%03d.jpg`
	cmd := fmt.Sprintf(GetFrameCmd, src, dst, f.String())
	_, err := osi.CMD(cmd)
	return err
}
