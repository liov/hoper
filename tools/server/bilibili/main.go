package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/tool"

	"tools/bilibili/parser"
)

func main() {
	var page int
	var continuous bool
	flag.IntVar(&page, "p", 1, "收藏页数")
	flag.StringVar(&tool.DownloadPath, "d", "D:/F/B站", "收藏页数")
	flag.BoolVar(&continuous, "l", false, "是否连续")
	flag.Parse()

	aid := tool.Bv2av("BV1Br4y1j7pa")
	req := parser.GetRequestByFav(aid)
	crawler.New(req)

	if continuous {
		req := parser.FavReq(page)
		crawler.New(req)
	} else {
		req := parser.FavReqs(page)
		crawler.New(req...)
	}
}
