package main

import (
	"time"

	"tools/pro"
)

func main() {
	pro.SetDB()
	pro.Start(normal)
}

func normal(sd *pro.Speed) {
	start := 370000
	end := 400000
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go pro.Fetch(i, sd)
		time.Sleep(pro.Interval)
	}
}

func one(id int) {
	/*	reader, _ := pro.Request(http.DefaultClient, fmt.Sprintf(pro.CommonUrl, "375509"))
		data, _ := ioutil.ReadAll(reader)
		log.Println(string(data))*/
	sd := pro.NewSpeed(pro.Loop)
	sd.WebAdd(1)
	pro.Fetch(id, sd)
	sd.Wait()
}
