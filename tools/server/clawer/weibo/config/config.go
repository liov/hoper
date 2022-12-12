package config

import (
	"github.com/liov/hoper/server/go/lib/utils/conctrl"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"path/filepath"
	"tools/clawer/weibo/rpc"
)

type Customize struct {
	UserId            int
	StopTime          int
	WorkCount         uint
	DownloadPath      string
	DownloadVideoPath string
	DownloadPicPath   string
	DownloadTmpPath   string
	Cookie            string
	SkipKind          []conctrl.Kind
	Users             []int
}

type config struct {
	Weibo Customize
	//Log      log.LogConfig
}

func (c *config) Init() {
	if c.Weibo.WorkCount == 0 {
		c.Weibo.WorkCount = 5
	}
	if c.Weibo.StopTime == 0 {
		c.Weibo.StopTime = 1
	}
	rpc.Cookie = c.Weibo.Cookie
	c.Weibo.DownloadPath, _ = filepath.Abs(c.Weibo.DownloadPath + "/debug1")
	c.Weibo.DownloadVideoPath = c.Weibo.DownloadPath + fs.PathSeparator + "video"
	c.Weibo.DownloadPicPath = c.Weibo.DownloadPath + fs.PathSeparator + "pic"
	c.Weibo.DownloadTmpPath = c.Weibo.DownloadPath + fs.PathSeparator + "tmp"

}

var Conf = &config{}
