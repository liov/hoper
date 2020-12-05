package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/liov/hoper/go/v2/utils/fs"
	py "tools/pinyin"
	"tools/pro"
)

func main() {
	pro.SetDB()
	pro.Start(history)
}

func history(sd *pro.Speed) {
	start := 353182
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
		invalidPost := &pro.InvalidPost{TId: id, Reason: 0}
		err := pro.DB.Save(invalidPost).Error
		if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			sd.FailInvalidPost <- tid
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
	if s.Length() < 1 {
		invalidPost := &pro.InvalidPost{TId: id, Reason: 1}
		err = pro.DB.Save(invalidPost).Error
		if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			sd.FailInvalidPost <- tid + " 0"
		}
		return
	}
	auth, title, postTime, post := parseHtmlHistory(doc)
	post.TId = id
	err = pro.DB.Save(post).Error
	if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
		sd.FailPost <- tid + " 1"
	}

	dir := pro.CommonDir

	if auth != "" {
		dir += py.FistLetter(auth) + pro.Sep + auth + pro.Sep
	}
	if title != "" {
		dir += title + `_` + tid + pro.Sep
	}

	dir = fs.PathClean(dir)
	_, err = os.Stat(dir + `content.txt`)
	if !os.IsNotExist(err) {
		os.Rename(dir+`content.txt`, dir+postTime+`.txt`)
		log.Println("rename:", dir+postTime)
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
