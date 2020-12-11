package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"tools/pro"
)

func main() {
	pro.SetDB()
	//one(351556)
	//test(401100)
	pro.Start(normal)
}

func normal(sd *pro.Speed) {
	start := 266110
	end := 300000
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go pro.Fetch(i, sd)
		time.Sleep(pro.Interval)
	}
}

func one(id int) {
	sd := pro.NewSpeed(pro.Loop)
	sd.WebAdd(1)
	pro.Fetch(id, sd)
	sd.Wait()
}

func test(id int) {
	reader, err := pro.Request(http.DefaultClient, fmt.Sprintf(pro.CommonUrl, strconv.Itoa(id)))
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(reader)
	log.Println(string(data))
}
