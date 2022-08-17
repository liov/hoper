package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"time"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	req := download.GetViewInfoReq(984161801, download.CoverViewInfoHandleFun)
	crawler.New(10).SkipKind(4).Timer(1, time.Millisecond*500).Timer(3, time.Millisecond*500).Run(req)
}
