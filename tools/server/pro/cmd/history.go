package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"tools/pro"
)

func main() {
	pro.SetDB()
	pro.Start(history)
}

func history(sd *pro.Speed) {
	start := 300100
	end := 400000
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetchHistory(i, sd)
		time.Sleep(pro.Interval)
	}
}

func fetchHistory(id int, sd *pro.Speed) {
	defer sd.WebDone()
	tid := strconv.Itoa(id)
	reader, err := pro.Request(http.DefaultClient, fmt.Sprintf(pro.CommonUrl, tid))
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
	post.PicNum = int8(s.Length())
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
		post.CreatedAt = postTime[len(`发表于 `):]
	}
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
