package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/config"
	"tools/bilibili/download"
)

func main() {

	/*	aid := tool.Bv2av("BV1Br4y1j7pa")
		req := parser.GetViewInfoReq(aid)
		crawler.New(req)
	*/

	defer initialize.Start(config.Conf, nil)()
	req := download.FavReqs(config.Conf.Bilibili.PageStart, config.Conf.Bilibili.PageEnd)
	crawler.New(10).Run(req...)
	/*	req := download.GetByBvId("BV1AB4y187HK")
		crawler.New(10).Run(req)*/
}
