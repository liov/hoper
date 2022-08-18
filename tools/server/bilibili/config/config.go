package config

import (
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"os"
	"path/filepath"
	"tools/bilibili/rpc"
)

type Customize struct {
	PageBegin         int `init:flag:p`
	PageEnd           int
	DownloadPath      string
	DownloadVideoPath string
	DownloadPicPath   string
	Cookie            string
	SkipKind          []conctrl.Kind
}

type config struct {
	Bilibili Customize
	//Log      log.LogConfig
}

func (c *config) Init() {
	if c.Bilibili.PageEnd == 0 {
		c.Bilibili.PageEnd = 1
	}
	if c.Bilibili.PageBegin == 0 {
		c.Bilibili.PageBegin = c.Bilibili.PageEnd
	}
	rpc.Cookie = c.Bilibili.Cookie
	c.Bilibili.DownloadPath, _ = filepath.Abs(c.Bilibili.DownloadPath)
	c.Bilibili.DownloadVideoPath = c.Bilibili.DownloadPath + "/video"
	err := os.MkdirAll(c.Bilibili.DownloadVideoPath, 0777)
	if err != nil {
		panic(err)
	}
	c.Bilibili.DownloadPicPath = c.Bilibili.DownloadPath + "/pic"
	err = os.MkdirAll(c.Bilibili.DownloadPicPath, 0777)
	if err != nil {
		panic(err)
	}
}

var Conf = &config{}
