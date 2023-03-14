package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/v2/utils/net/http/client/crawler"
	"os"
	"time"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/download"
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
	engine := crawler.NewEngine(config.Conf.Weibo.WorkCount).Timer(download.KindGet, time.Second)

	engine.Run(download.GetUserAllFollowsReq(config.Conf.Weibo.UserId))
}
