package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	wg := new(sync.WaitGroup)
	url := "https://f1113.wonderfulday27.live/viewthread.php?tid=368996"
	reader, err := request(url)
	if err != nil {
		log.Println(err)
	}
	defer reader.Close()
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(`img[src="images/common/none.gif"]`).Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("file"); ok {
			wg.Add(1)
			go download(url, wg)
		}
	})
	wg.Wait()
}

func download(url string, wg *sync.WaitGroup) {
	reader, err := request(url)
	if err != nil {
		log.Println(err)
	}
	defer reader.Close()
	f, _ := os.Create(`F:\pic\` + strings.Split(url, "//")[2])
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		log.Println(err)
	}
	log.Println("下载成功：", url)
	wg.Done()
}

func request(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	http.DefaultClient.Timeout = 300 * time.Second
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		reader = resp.Body
	}
	return reader, nil
}
