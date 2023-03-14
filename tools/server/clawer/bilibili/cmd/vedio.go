package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/v2/utils/net/http/client/crawler"
	"os"

	"time"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	err := os.MkdirAll(config.Conf.Bilibili.DownloadTmpPath, 0777)
	if err != nil {
		panic(err)
	}
	engine := crawler.NewEngine(config.Conf.Bilibili.WorkCount).StopAfter(time.Hour * time.Duration(config.Conf.Bilibili.StopTime))
	go download.DownloadRecordVideo(engine)
	download.GetEngineMerge(engine)
	time.Sleep(time.Second)
	engine.Run()
}
