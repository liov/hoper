package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/v2/utils/net/http/client/crawler"
	"time"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	engine := crawler.NewEngine(config.Conf.Weibo.WorkCount).Timer(download.KindGet, time.Second)
	engine.Run(download.RecordUsersWeiboReq(config.Conf.Weibo.Users, true)...)
}
