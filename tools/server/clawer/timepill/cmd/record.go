package main

import (
	"context"
	"flag"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/conctrl/listener"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"time"
	"tools/clawer/timepill"
)

var today = flag.Bool("t", false, "记录今天日记")

// go build -o timepill/timepill timepill/cmd/record.go
// 日记记录
func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()

	//go timepill.RecordByOrderUser()
	flag.Parse()
	if *today {
		log.Info("todayRecord")
		timepill.TodayRecord()
	}
	task := &listener.TimerTask{}
	task.Do = func(ctx context.Context) {
		log.Info("定时任务：记录评论执行,times:", task.Times)
		timepill.CronCommentRecord()
	}
	go task.Timer(context.Background(), time.Hour)

	//go timepill.RecordByOrderNoteBook()
	recordtask := &listener.TimerTask{}
	recordtask.Do = func(ctx context.Context) {
		log.Info("定时任务：记录评论执行,times:", task.Times)
		timepill.RecordTask()
	}
	recordtask.Timer(context.Background(), time.Second*timepill.Conf.TimePill.Timer)
}
