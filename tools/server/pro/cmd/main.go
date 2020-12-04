package main

import (
	"strconv"
	"time"

	"tools/pro"
)

func main() {
	pro.Start(normal)
}

func normal(sd *pro.Speed) {
	start := 370000
	end := 400000
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go pro.Fetch(strconv.Itoa(i), sd)
		time.Sleep(pro.Interval)
	}
}
