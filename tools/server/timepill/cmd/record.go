package main

import (
	"flag"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"tools/timepill"
)

var today = flag.Bool("t", false, "记录今天日记")

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()

	/*	go timepill.RecordByOrderUser()
		log.Info("startRecord")
		if *today {
			log.Info("todayRecord")
			timepill.TodayRecord()
		}
		log.Info("startRecord")*/
	timepill.StartRecord()
}
