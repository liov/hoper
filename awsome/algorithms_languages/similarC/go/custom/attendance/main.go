package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/robfig/cron"
	"github.com/tidwall/gjson"
)

var kqReq *http.Request
var lastId int64
var ding Ding
var at = map[string]string{
	"000002204": "xxx",
	"000002190": "xxx",
}

type Ding struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

type Text struct {
	Content string `json:"content"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func main() {
	log.SetFlags(15)
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch,
		// kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	urlStr := `http://218.17.157.34:1234/grid/att/CheckInOutGrid/`
	kqReq, _ = http.NewRequest("POST", urlStr, strings.NewReader("page=1&rp=10"))
	kqReq.Header.Set("Cookie", "sessionidadms=xxx")
	kqReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	Request()
	c := cron.New()
	c.AddFunc("*/20 * 0,8,9,13,18,19,20,21,22,23 * * *", Request)
	c.Start()
	<-ch
}

func Request() {
	resp, err := http.DefaultClient.Do(kqReq)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	array := gjson.Get(string(body), "rows").Array()

	for i := len(array) - 1; i >= 0; i-- {
		obj := array[i].Map()
		id := obj["id"].Int()
		name := obj["name"].String()
		if utf8.RuneCountInString(name) == 2 {
			name = name + "    "
		}
		depName := obj["DeptName"].String()
		checktime := obj["checktime"].String()

		if id > lastId && depName == "平台研发中心" {
			lastId = id
			ding.MsgType = "text"
			ding.Text.Content = ding.Text.Content + name + ` : ` + checktime + "\n"
			for k, v := range at {
				if obj["badgenumber"].String() == k {
					ding.At.AtMobiles = append(ding.At.AtMobiles, v)
				}
			}

		}
	}
	if ding.Text.Content == "" {
		return
	}

	ding.Text.Content = ding.Text.Content[:len(ding.Text.Content)-1]
	body, _ = json.Marshal(&ding)
	urlStr := `https://oapi.dingtalk.com/robot/send?access_token=xxx`
	dingReq, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	dingReq.Header.Set("Content-Type", "application/json")
	log.Println("请求钉钉")
	dresp, err := http.DefaultClient.Do(dingReq)
	if err != nil {
		log.Println(err)
		return
	}
	ding.Text.Content = ""
	ding.At.AtMobiles = ding.At.AtMobiles[:0]
	defer resp.Body.Close()
	dbody, err := ioutil.ReadAll(dresp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(dbody))

}
