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
	c := cron.New()
	c.AddFunc("0 55 23 * * ?", timepill.TodayCommentRecord)
	//go timepill.RecordByOrderNoteBook()
	timepill.StartRecord()
}