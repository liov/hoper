package main

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const commonUrl = "https://f1113.wonderfulday27.live/viewthread.php?tid=%d"
const loop = 50
const commonDir = `E:\pic\`

var userAgent = []string{
	`Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36`,
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36",
	"Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10.5; en-US; rv:1.9.2.15) Gecko/20110303 Firefox/3.6.15",
	`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36`,
}

var reqCache = make([]*http.Request, 1, loop)
var picClient = new(http.Client)

func init() {
	req, _ := newRequest(commonUrl)
	reqCache[0] = req.Clone(context.Background())
	http.DefaultClient.Timeout = 300 * time.Second
	picClient.Timeout = 300 * time.Second
	proxyURL, _ := url.Parse("socks5://localhost:8080")
	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
	}
	picClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
	}
}

func main() {
	start := 360000
	end := 370000
	sd := NewSpeed(loop)
	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)
		f, _ := os.Create(commonDir + "fail_" + time.Now().Format("2006_01_02_15_04_05") + `.txt`)
		for txt := range sd.fail {
			f.WriteString(txt + "\n")
		}
		f.Close()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		f, _ := os.Create(commonDir + "fail_pic_" + time.Now().Format("2006_01_02_15_04_05") + `.txt`)
		for txt := range sd.failPic {
			f.WriteString(txt + "\n")
		}
		f.Close()
		wg.Done()
	}()
	for i := start; i < end; i++ {
		sd.Add(1)
		go fetch(i, sd)
		time.Sleep(time.Second)
	}
	sd.Wait()
	close(sd.fail)
	close(sd.failPic)
	wg.Wait()
}

type speed struct {
	wg            *sync.WaitGroup
	c             chan struct{}
	fail, failPic chan string
}

func (s *speed) Add(i int) {
	s.wg.Add(i)
	s.c <- struct{}{}
}

func (s *speed) Done() {
	s.wg.Done()
	<-s.c
}

func (s *speed) Wait() {
	s.wg.Wait()
}

func NewSpeed(cap int) *speed {
	return &speed{
		wg:      new(sync.WaitGroup),
		c:       make(chan struct{}, cap),
		fail:    make(chan string, cap),
		failPic: make(chan string, cap),
	}
}

func fetch(id int, wg *speed) {
	defer wg.Done()
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
		log.Fatal(err)
	}
	reader.Close()
	s := doc.Find(`img[src="images/common/none.gif"]`)
	if s.Length() < 1 {
		return
	}
	auth, title, text, htl := parseHtml(doc)
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
	if text != "" {
		f, err := os.Create(dir + `content.txt`)
		f.WriteString(text)
		f.Close()
		if err != nil {
			log.Println(err)
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

	s.Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("file"); ok {
			wg.Add(1)
			go download(url, dir, wg)
			time.Sleep(time.Second)
		}
	})
}

func parseHtml(doc *goquery.Document) (string, string, string, *goquery.Selection) {
	auth := doc.Find(".mainbox td.postauthor .postinfo a").First().Text()
	title := doc.Find("#threadtitle h1").Text()
	content := doc.Find(".t_msgfont").First()
	text := content.Contents().Not(".t_attach").Text()
	html := content.Not(".t_attach").Not("span")
	auth = strings.ReplaceAll(auth, " ", "")
	title = strings.ReplaceAll(title, " ", "")
	return auth, title, text, html
}

func download(url, dir string, wg *speed) {
	defer wg.Done()

	reader, err := request(picClient, url)
	if err != nil {
		log.Println(err, "url:", url)
		if !strings.HasPrefix(err.Error(), "返回错误") {
			wg.failPic <- url + " " + dir
		}
		return
	}
	defer reader.Close()
	s := strings.Split(url, "//")
	name := s[len(s)-1]
	if strings.Contains(name, "/") {
		s = strings.Split(url, "/")
		name = s[len(s)-1]
	}
	if strings.Contains(name, "\\") {
		s = strings.Split(url, "\\")
		name = s[len(s)-1]
	}
	f, err := os.Create(dir + name)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		log.Println("写入文件错误：", err)
	}
	log.Printf("下载成功：%s,目录：%s\n", url, dir)
}

func request(client *http.Client, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	var reader io.ReadCloser
	var resp *http.Response
	for i := 0; i < 20; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		n := rand.Intn(5)
		req.Header.Set("User-Agent", userAgent[n])
		resp, err = client.Do(req)
		if err != nil {
			log.Println(err, "url:", url)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("返回错误，状态码：%d,url:%s", resp.StatusCode, url)
		}

		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				if resp != nil {
					resp.Body.Close()
				}
				log.Println(err, "url:", url)
				continue
			}
		} else {
			reader = resp.Body
		}
		if reader != nil {
			break
		}
	}
	if reader == nil {
		if resp != nil {
			resp.Body.Close()
		}
		msg := "请求失败：" + url
		if err != nil {
			msg = err.Error() + msg
		}
		return nil, errors.New(msg)
	}
	return reader, nil
}

func newRequest(url string) (*http.Request, error) {
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
	return req, nil
}
