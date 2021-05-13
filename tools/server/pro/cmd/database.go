package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/liov/hoper/v2/tools/create-table/get"

	"tools/pro"
)

func main() {
	pro.SetDB()
	pro.Start(history)
}

func history(sd *pro.Speed) {
	start := 407889
	end := 407912
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetchHistory(i, sd)
		time.Sleep(pro.Interval)
	}
}

func database() {
	get.GetDB().Migrator().CreateTable(&pro.Post{})
}

func historyOne(sd *pro.Speed) {
	historyFormFile(`fail_post_2020_12_08_13_07_15`, sd)
}

func historyFormFile(path string, sd *pro.Speed) {
	f, err := os.Open(pro.CommonDir + path + pro.Ext)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		tid, status := s[0], s[1]
		id, _ := strconv.Atoi(tid)
		if status == "2" {
			invalidPost := &pro.Post{TId: id, Status: 2}
			err := pro.DB.Save(invalidPost).Error
			if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
				sd.FailDB <- tid + " 2"
			}
			continue
		}
		sd.WebAdd(1)
		go fetchHistory(id, sd)
		time.Sleep(pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func fetchHistory(id int, sd *pro.Speed) {
	defer sd.WebDone()
	tid := strconv.Itoa(id)
	reader, err := pro.Request(http.DefaultClient, pro.CommonUrl+tid)
	if err != nil {
		//log.Println(err, "id:", tid)
		if !strings.HasPrefix(err.Error(), "返回错误") {
			sd.Fail <- tid
		}
		invalidPost := &pro.Post{TId: id, Status: 2}
		err := pro.DB.Save(invalidPost).Error
		if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			sd.FailDB <- tid + " 2"
		}
		return
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
		sd.Fail <- tid
		return
	}
	reader.Close()
	s := doc.Find(`img[src="images/common/none.gif"]`)
	_, _, _, post := parseHtmlHistory(doc)
	post.TId = id
	post.PicNum = uint32(s.Length())
	status := "0"
	if post.PicNum == 0 {
		status = "1"
	}
	err = pro.DB.Save(post).Error
	if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
		sd.FailDB <- tid + " " + status
	}
}

func parseHtmlHistory(doc *goquery.Document) (string, string, string, *pro.Post) {
	auth := doc.Find(".mainbox td.postauthor .postinfo a").First().Text()
	title := doc.Find("#threadtitle h1").Text()
	postTime := doc.Find(".authorinfo em").First().Text()
	post := &pro.Post{
		TId:   0,
		Auth:  auth,
		Title: title,
	}
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

	post.CreatedAt = postTime
	postTime = strings.ReplaceAll(postTime, ":", "-")

	auth = strings.ReplaceAll(auth, "\\", "")
	auth = strings.ReplaceAll(auth, "/", "")
	auth = strings.ReplaceAll(auth, ":", "")
	if strings.HasSuffix(auth, ".") {
		auth += "$"
	}
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, ":", "")
	if strings.HasSuffix(title, ".") {
		title += "$"
	}

	postTime = strings.ReplaceAll(postTime, ":", "-")
	content := doc.Find(".t_msgfont").First()
	text := content.Contents().Not(".t_attach").Text()
	post.Content = text
	return auth, title, postTime, post
}
