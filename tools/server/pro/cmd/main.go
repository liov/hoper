package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"time"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()

	pro.Start(normal)
}

func normal(sd *pro.Speed) {
	start := 492010
	end := 492020
	for i := start; i <= end; i++ {
		sd.WebAdd(1)
		go pro.Fetch(i, sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
}
