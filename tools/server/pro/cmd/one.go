package main

import (
	"time"

	"tools/pro"
)

func main() {
	pro.SetDB()
	sd := pro.NewSpeed(pro.Loop)
	start := 211301
	end := 211302
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		pro.Fetch(i, sd)
		time.Sleep(pro.Interval)
	}
	sd.Wait()
}
