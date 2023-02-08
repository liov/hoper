package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/conctrl"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"log"
	"time"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()

	task1 := &conctrl.TimerTask{}
	task1.Do = func(ctx context.Context) {
		log.Println("task1 times", task1.Times)
		engine := crawler.NewEngine(config.Conf.Weibo.WorkCount).Timer(download.KindGet, time.Second)

		engine.Run(download.GetUserFollowWeiboReq(""))
	}
	go task1.Timer(context.Background(), time.Minute)
	engine := crawler.NewEngine(config.Conf.Weibo.WorkCount).Timer(download.KindGet, time.Second)
	engine.Run(download.RecordUsersWeiboReq(config.Conf.Weibo.Users, true)...)
	task := &conctrl.TimerTask{}
	task.Do = func(ctx context.Context) {
		log.Println("task times", task.Times)
		engine.ReRun(download.RecordUsersWeiboReq(config.Conf.Weibo.Users, true)...)
	}
	task.RandTimer(context.Background(), time.Minute*10, time.Minute*20)
}
