package config

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/log"
	"os"
	"tools/bilibili/rpc"
)

type Customize struct {
	PageStart       int
	PageEnd         int
	DownloadPath    string
	DownloadPicPath string
	Cookie          string
}

type config struct {
	Bilibili Customize
	Log      log.LogConfig
}

func (c *config) Init() {
	if c.Bilibili.PageEnd == 0 {
		c.Bilibili.PageEnd = 1
	}
	if c.Bilibili.PageStart == 0 {
		c.Bilibili.PageStart = c.Bilibili.PageEnd
	}
	if c.Bilibili.PageStart > c.Bilibili.PageEnd {
		c.Bilibili.PageEnd = c.Bilibili.PageStart
	}
	rpc.Cookie = c.Bilibili.Cookie
	err := os.MkdirAll(c.Bilibili.DownloadPath, 0777)
	if err != nil {
		panic(err)
	}
	c.Bilibili.DownloadPicPath = c.Bilibili.DownloadPath + "_pic"
	err = os.MkdirAll(c.Bilibili.DownloadPicPath, 0777)
	if err != nil {
		panic(err)
	}
}

var Conf = &config{}
