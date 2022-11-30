package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"os"
	"time"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	err := os.MkdirAll(config.Conf.Bilibili.DownloadVideoPath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(config.Conf.Bilibili.DownloadTmpPath, 0777)
	if err != nil {
		panic(err)
	}
	engine := crawler.NewEngine(config.Conf.Bilibili.WorkCount).Timer(download.KindViewInfo, time.Second).SkipKind(download.KindDownloadVideo).Timer(download.KindGetPlayerUrl, time.Second).StopAfter(time.Hour * time.Duration(config.Conf.Bilibili.StopTime))
	go download.FixRecordFav(engine)

	engine.RunSingleWorker()
}
