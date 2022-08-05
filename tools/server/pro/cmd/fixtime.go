package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"log"
	"strconv"
	"strings"
	"time"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	//test(401100)
	pro.Start(fixtime)
}

func fixtime(sd *pro.Speed) {
	start := 456610
	end := 458600
	for i := start; i <= end; i++ {
		sd.WebAdd(1)
		go doFixtime(i, sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
}

func doFixtime(id int, sd *pro.Speed) {
	defer sd.WebDone()
	tid := strconv.Itoa(id)
	reader, err := pro.R(pro.Conf.Pro.CommonUrl + tid)
	if err != nil {
		//log.Println(err, "id:", tid)
		return
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
		sd.Fail <- tid
		return
	}

	postTime := doc.Find(".posterinfo .authorinfo em").First().Text()
	if strings.HasPrefix(postTime, "发表于") {
		postTime = postTime[len(`发表于 `):]
	}
	if strings.Contains(postTime, "天") {
		now := time.Now()
		var day int
		if strings.Contains(postTime, "天前") {
			day, _ = strconv.Atoi(postTime[0:0])
		} else {
			describe := postTime[:6]
			switch describe {
			case "前天":
				day = 2
			case "昨天":
				day = 1
			case "今天":
				day = 0
			}
		}
		now.AddDate(0, 0, -day)
		date := now.Format("2006-01-02")
		postTime = date + " " + postTime[len(postTime)-5:]
	}
	pro.Dao.DB.Exec(`Update post SET created_at = ? Where t_id = ?`, postTime, id)
}
