package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/robfig/cron/v3"
	"time"
	"tools/timepill"
)

var today = flag.Bool("t", false, "记录今天日记")

// go build -o timepill/timepill timepill/cmd/record.go
// 日记记录
func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()

	//go timepill.RecordByOrderUser()

	if *today {
		log.Info("todayRecord")
		timepill.TodayRecord()
	}
	c := cron.New()
	_, err := c.AddFunc("0 */1 * * *", func() {
		log.Info("定时任务：记录评论执行")
		timepill.CronCommentRecord()
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	//go timepill.RecordByOrderNoteBook()
	timepill.RecordTask()
	tc := time.NewTicker(time.Second * timepill.Conf.TimePill.Timer)
	for range tc.C {
		timepill.RecordTask()
	}
}
