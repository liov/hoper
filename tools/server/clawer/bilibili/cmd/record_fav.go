package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/conctrl/listener"
	"github.com/liov/hoper/server/go/lib/v2/utils/net/http/client/crawler"
	"log"
	"time"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/download"

	"tools/clawer/bilibili/config"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	task := &listener.TimerTask{}
	engine := crawler.NewEngine(config.Conf.Bilibili.WorkCount).SkipKind(download.KindDownloadVideo).Timer(download.KindViewInfo, time.Second).Timer(download.KindGetPlayerUrl, time.Second)
	engine.SpeedLimited(time.Second)
	task.Do = func(ctx context.Context) {
		log.Println("times", task.Times)
		/*req1 := download.FavReqs(63181530, 1, 5, download.GetFavList)
		req2 := download.FavReqs(62504730, 1, 1, download.GetFavList)
		req := append(req1, req2...)*/
		engine.ReRun(download.RecordFavTimer(time.Now())...)
	}
	task.FirstExec = true
	task.Timer(context.Background(), time.Minute)
}
