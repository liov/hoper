package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	timepill.RecordByOrderUser()
}
