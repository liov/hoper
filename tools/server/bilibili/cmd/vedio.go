package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/initialize"
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
	engine := crawler.New(config.Conf.Bilibili.WorkCount).StopAfter(time.Hour * time.Duration(config.Conf.Bilibili.StopTime))
	go download.FavVideo(engine)
	engine.Run()
}
