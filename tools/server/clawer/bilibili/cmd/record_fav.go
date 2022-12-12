package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/conctrl"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"log"
	"time"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/download"

	"tools/clawer/bilibili/config"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	task := &conctrl.TimerTask{}
	task.Do = func(ctx context.Context) {
		log.Println("times", task.Times)
		/*req1 := download.FavReqs(63181530, 1, 5, download.GetFavList)
		req2 := download.FavReqs(62504730, 1, 1, download.GetFavList)
		req := append(req1, req2...)*/
		engine := crawler.NewEngine(config.Conf.Bilibili.WorkCount).SkipKind(download.KindDownloadVideo).Timer(download.KindViewInfo, time.Second).Timer(download.KindGetPlayerUrl, time.Second)
		engine.SpeedLimited(time.Second)
		engine.Run(download.RecordFavTimer(time.Now())...)
	}
	task.RandTimer(context.Background(), time.Minute*10, time.Minute*20)
}
