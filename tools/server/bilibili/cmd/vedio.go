package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"time"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/download"
)

func main() {
	flag.IntVar(&config.Conf.Bilibili.PageBegin, "pb", 1, "开始页")
	flag.IntVar(&config.Conf.Bilibili.PageEnd, "pe", 1, "开始页")
	defer initialize.Start(config.Conf, &dao.Dao)()
	flag.Parse()
	req := download.FavReqs(config.Conf.Bilibili.PageBegin, config.Conf.Bilibili.PageEnd, download.FavList)
	crawler.New(config.Conf.Bilibili.WorkCount).Timer(1, time.Millisecond*500).Timer(3, time.Second*10).Run(req...)
}
