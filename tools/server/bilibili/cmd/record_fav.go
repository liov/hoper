package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"log"
	"time"
	"tools/bilibili/dao"

	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/config"
	"tools/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	task := &conctrl.TimerTask{}
	task.Do = func(ctx context.Context) {
		log.Println("times", task.Times)
		req := download.FavReqs(1, 1, download.RecordFavList)
		crawler.New(10).SkipKind(4).Timer(1, time.Millisecond*500).Timer(3, time.Millisecond*500).Run(req...)
	}
	conctrl.Timer(context.Background(), task, time.Minute*30)
}
