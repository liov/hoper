package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"time"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()

	pro.Start(normal)
}

func normal(sd *pro.Speed) {
	start := 510572
	end := 510683
	for i := start; i <= end; i++ {
		sd.WebAdd(1)
		go pro.Fetch(i, sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
}
