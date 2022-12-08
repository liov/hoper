package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"os"
	"time"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/download"

	"tools/clawer/weibo/config"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	err := os.MkdirAll(config.Conf.Weibo.DownloadVideoPath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(config.Conf.Weibo.DownloadTmpPath, 0777)
	if err != nil {
		panic(err)
	}
	engine := crawler.NewEngine(config.Conf.Weibo.WorkCount).Timer(download.KindGetFollow, time.Second).Timer(download.KindGetPhoto, time.Second)

	engine.Run(download.GetUserAllFollowsReq(config.Conf.Weibo.UserId))
}
