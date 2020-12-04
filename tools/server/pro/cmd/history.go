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
	"golang.org/x/net/html"
	py "tools/pinyin"
	"tools/pro"
)

func main() {
	pro.Start(history)
}

func history(sd *pro.Speed) {
	start := 300000
	end := 360000
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetchHistory(strconv.Itoa(i), sd)
		time.Sleep(pro.Interval)
	}
}

func fetchHistory(tid string, sd *pro.Speed) {
	defer sd.WebDone()
	reader, err := pro.Request(http.DefaultClient, fmt.Sprintf(pro.CommonUrl, tid))
	if err != nil {
		log.Println(err, "id:", tid)
		if !strings.HasPrefix(err.Error(), "返回错误") {
			sd.Fail <- tid
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
		return
	}
	auth, title, text, htl, postTime, his := parseHtmlHistory(doc)

	dir := pro.CommonDir

	if auth != "" {
		dir += py.FistLetter(auth) + pro.Sep + auth + pro.Sep
	}
	if title != "" {
		dir += title + `_` + tid + pro.Sep
	}
	if his {
		dir = fs.PathClean(dir)
	}

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			log.Println(err, dir)
			sd.Fail <- tid
			return
		}
	}
	id, _ := strconv.Atoi(tid)
	if his && text != "" {
		f, err := os.Create(dir + postTime + `.txt`)
		f.WriteString(text)
		f.Close()
		if err != nil {
			log.Println(err)
		}
		if id < 307000 {
			os.Remove(dir + auth + `.txt`)
		}
		if htl.Length() > 0 {
			f, err = os.Create(dir + `index.html`)
			for c := htl.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
				html.Render(f, c)
			}
			f.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}
	if !his {

		if id < 307000 {
			os.Rename(dir+auth+`.txt`, dir+postTime+`.txt`)
		} else {
			os.Rename(dir+`content.txt`, dir+postTime+`.txt`)
		}

		return
	}
	s.Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("file"); ok {
			sd.Add(1)
			go pro.Download(url, dir, sd)
			time.Sleep(pro.Interval)

		}
	})
}

func parseHtmlHistory(doc *goquery.Document) (string, string, string, *goquery.Selection, string, bool) {
	auth := doc.Find(".mainbox td.postauthor .postinfo a").First().Text()
	postTime := doc.Find(".authorinfo em").First().Text()
	postTime = strings.ReplaceAll(postTime, ":", "-")
	title := doc.Find("#threadtitle h1").Text()
	content := doc.Find(".t_msgfont").First()

	if strings.Contains(auth, "\\") ||
		strings.Contains(auth, "/") ||
		strings.Contains(auth, ":") ||
		strings.Contains(auth, "<") ||
		strings.Contains(auth, ">") ||
		strings.Contains(auth, "\"") ||
		strings.Contains(auth, "|") ||
		strings.Contains(auth, "?") ||
		strings.Contains(auth, "*") ||
		strings.Contains(title, "\\") ||
		strings.Contains(title, "/") ||
		strings.Contains(title, ":") ||
		strings.Contains(title, "<") ||
		strings.Contains(title, ">") ||
		strings.Contains(title, "\"") ||
		strings.Contains(title, "|") ||
		strings.Contains(title, "?") ||
		strings.Contains(title, "*") {
		auth = strings.ReplaceAll(auth, "\\", "")
		auth = strings.ReplaceAll(auth, "/", "")
		title = strings.ReplaceAll(title, "\\", "")
		title = strings.ReplaceAll(title, "/", "")
		text := content.Contents().Not(".t_attach").Text()
		html := content.Not(".t_attach").Not("span")
		return auth, title, text, html, postTime, true
	}
	return auth, title, "", nil, postTime, false
}
