package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func gif(start, end int, sd *speed) {
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetchGif(i, sd)
		time.Sleep(interval)
	}
}
func fetchGif(id int, wg *speed) {
	defer wg.WebDone()
	tid := strconv.Itoa(id)
	reader, err := request(http.DefaultClient, fmt.Sprintf(commonUrl, id))
	if err != nil {
		log.Println(err, "id:", tid)
		if !strings.HasPrefix(err.Error(), "返回错误") {
			wg.fail <- tid
		}
		return
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
		wg.fail <- tid
		return
	}
	reader.Close()
	s := doc.Find(`img[src="images/common/none.gif"]`)
	if s.Length() < 1 {
		return
	}
	auth, title := parseHtmlGif(doc)
	dir := commonDir
	if auth != "" {
		dir += auth + `\`
	}
	if title != "" {
		dir += title + `_` + tid + `\`
	}

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
				wg.Add(1)
				go download(url, dir, wg)
				time.Sleep(interval)
			}
		}
	})
}

func parseHtmlGif(doc *goquery.Document) (string, string) {
	auth := doc.Find(".mainbox td.postauthor .postinfo a").First().Text()
	title := doc.Find("#threadtitle h1").Text()
	auth = strings.ReplaceAll(auth, " ", "")
	title = strings.ReplaceAll(title, " ", "")
	return auth, title
}
