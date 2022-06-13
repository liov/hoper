package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	sd := pro.NewSpeed(pro.Conf.Pro.Loop)

	s := []int{
		434657,
	}
	for i := 0; i < len(s); i++ {
		sd.WebAdd(1)
		pro.Fetch(s[i], sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
	sd.Wait()
}

func test(id int) {
	reader, err := pro.Request(http.DefaultClient, pro.Conf.Pro.CommonUrl+strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(reader)
	log.Println(string(data))
}
