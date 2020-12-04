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
	pro.Start(gif)
}

func gif(sd *pro.Speed) {
	start := 370000
	end := 400000
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetchGif(strconv.Itoa(i), sd)
		time.Sleep(pro.Interval)
	}
}
func fetchGif(tid string, sd *pro.Speed) {
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
	auth, title := parseHtmlGif(doc)
	dir := pro.CommonDir

	if auth != "" {
		dir += py.FistLetter(auth) + pro.Sep + auth + pro.Sep
	}
	if title != "" {
		dir += title + `_` + tid + pro.Sep
	}
	dir = fs.PathClean(dir)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			log.Println(err, dir)
		}
	}

	s.Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("file"); ok {
			if strings.HasSuffix(url, ".gif") {
				sd.Add(1)
				go pro.Download(url, dir, sd)
				time.Sleep(pro.Interval)
			}
		}
	})
}

func parseHtmlGif(doc *goquery.Document) (string, string) {
	auth := doc.Find(".mainbox td.postauthor .postinfo a").First().Text()
	title := doc.Find("#threadtitle h1").Text()
	auth = strings.ReplaceAll(auth, "\\", "")
	auth = strings.ReplaceAll(auth, "/", "")
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, "/", "")
	return auth, title
}
