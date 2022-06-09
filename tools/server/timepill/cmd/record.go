package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/robfig/cron/v3"
	"tools/timepill"
)

var today = flag.Bool("t", false, "记录今天日记")

// go build -o timepill/timepill timepill/cmd/record.go
func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()

	//go timepill.RecordByOrderUser()

	if *today {
		log.Info("todayRecord")
		timepill.TodayRecord()
	}
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 */1 * * *", func() {
		log.Info("定时任务：记录评论执行")
		timepill.CronCommentRecord()
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	//go timepill.RecordByOrderNoteBook()
	timepill.StartRecord()
}
