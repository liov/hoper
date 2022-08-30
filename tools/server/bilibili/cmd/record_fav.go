package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"log"
	"time"
	"tools/bilibili/dao"
	"tools/bilibili/download"

	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/config"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	task := &conctrl.TimerTask{}
	task.Do = func(ctx context.Context) {
		log.Println("times", task.Times)
		/*req1 := download.FavReqs(63181530, 1, 5, download.RecordFavList)
		req2 := download.FavReqs(62504730, 1, 1, download.RecordFavList)
		req := append(req1, req2...)*/
		engine := crawler.New(config.Conf.Bilibili.WorkCount).SkipKind(4).Timer(1, time.Millisecond*500).Timer(3, time.Second)
		download.RecordFav(ctx, engine)
		engine.Run()
	}
	conctrl.Timer(context.Background(), task, time.Minute*10)
}
