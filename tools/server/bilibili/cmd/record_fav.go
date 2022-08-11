package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"time"
	"tools/bilibili/dao"

	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/config"
	"tools/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	req := download.FavReqs(config.Conf.Bilibili.PageStart, config.Conf.Bilibili.PageEnd)
	crawler.New(10).ExcludeKind(2, 3).Timer(1, time.Millisecond*500).Run(req...)
}
