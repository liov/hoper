package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/api"
	"tools/bilibili/download"
	"tools/bilibili/tool"
)

type Customize struct {
	Page         int
	DownloadPath string
	Continuous   bool
	Cookie       string
}
type config struct {
	Bilibili Customize
}

func (c config) Init() {
	api.Cookie = c.Bilibili.Cookie
	tool.DownloadPath = c.Bilibili.DownloadPath
}

func main() {

	/*	aid := tool.Bv2av("BV1Br4y1j7pa")
		req := parser.GetRequestByFav(aid)
		crawler.New(req)
	*/
	var c config
	defer initialize.Start(&c, nil)()
	if c.Bilibili.Continuous {
		req := download.FavReqs(c.Bilibili.Page)
		crawler.New(10).Run(req...)
	} else {
		req := download.FavReq(c.Bilibili.Page)
		crawler.New(10).Run(req)
	}
	/*	req := download.GetByBvId("BV1AB4y187HK")
		crawler.New(10).Run(req)*/
}
