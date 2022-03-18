package main

import (
	"time"

	"tools/pro"
)

func main() {
	pro.SetDB()
	//test(401100)
	pro.Start(normal)
}

func normal(sd *pro.Speed) {
	start := 476003
	end := 480600
	for i := start; i <= end; i++ {
		sd.WebAdd(1)
		go pro.Fetch(i, sd)
		time.Sleep(pro.Interval)
	}
}
