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
	sd := pro.NewSpeed(pro.Loop)

	s := []int{
		296314,
	}
	for i := 0; i < len(s); i++ {
		sd.WebAdd(1)
		pro.Fetch(s[i], sd)
		time.Sleep(pro.Interval)
	}
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
