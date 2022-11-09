package pro

import (
	"bufio"
	"bytes"
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	py2 "github.com/actliboy/hoper/server/go/lib/utils/strings/pinyin"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/strings"
	"golang.org/x/net/html"
)

var userAgent = []string{
	`Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36`,
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36",
	"Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10.5; en-US; rv:1.9.2.15) Gecko/20110303 Firefox/3.6.15",
	`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36`,
}

func GetFetchReq(id int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: "Fetch " + strconv.Itoa(id)}},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			s, dir, err := Fetch(id)
			if err != nil {
				return nil, err
			}
			if s != nil {
				var reqs []*crawler.Request
				s.Each(func(i int, s *goquery.Selection) {
					if url, ok := s.Attr("file"); ok {
						reqs = append(reqs, GetDownloadReq(url, dir))
					}
				})
				return reqs, nil
			}
			return nil, nil
		},
	}
}

func Fetch(id int) (*goquery.Selection, string, error) {

	tid := strconv.Itoa(id)
	reader, err := R(Conf.Pro.CommonUrl + tid)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "not found") {
			return nil, "", ReqPostError.Message(err.Error())
		}
		invalidPost := &Post{TId: id, Status: 2}
		err = Dao.DB.Save(invalidPost).Error
		if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			log.Println(err)
		}
		return nil, "", nil
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, "", ReqPostError.Message(err.Error())
	}
	s := doc.Find(`img[src="images/common/none.gif"]`)

	auth, title, text, postTime, htl, post := ParseHtml(doc)
	post.TId = id
	post.PicNum = uint32(s.Length())
	status := "0"
	if post.PicNum == 0 {
		post.Status = 1
		status = "1"
	}

	dir := Conf.Pro.CommonDir + "pic_" + strconv.Itoa(id/100000) + "/"

	if auth != "" {
		dir += py2.FistLetter(auth) + Sep + auth + Sep
	}
	if title != "" {
		dir += title + `_` + tid + Sep
	}
	dir = fs.PathClean(dir)

	post.Path = dir[CommonDirLen:]
	err = Dao.DB.Save(post).Error
	if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
		log.Println(err)
	}

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			log.Println(err, dir)
			return nil, "", MkdirError.Message(tid + " " + status).AppendErr(err)
		}
	}
	if text != "" {
		f, err := os.Create(dir + postTime + Conf.Pro.Ext)
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

	return s, dir, nil
}

func ParseHtml(doc *goquery.Document) (string, string, string, string, *goquery.Selection, *Post) {
	auth := doc.Find("#postlist .popuserinfo a").First().Text()
	title := doc.Find("#threadtitle h1").Text()
	timenode := doc.Find(".posterinfo .authorinfo em").First()
	postTime := timenode.Text()
	if strings.HasPrefix(postTime, "发表于") {
		postTime = postTime[len(`发表于 `):]
	}
	if strings.HasSuffix(postTime, "前") {
		postTime2, ok := timenode.Find("span").First().Attr("title")
		if !ok {
			postTime = time.Now().Format("2006-01-02 15:04:05")
		} else {
			postTime = postTime2
		}
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

	if strings.Contains(postTime, "刚才") {
		postTime = time.Now().Format("2006-01-02 15:04:05")
	}

	post := &Post{
		TId:   0,
		Auth:  auth,
		Title: title,
	}

	post.CreatedAt = postTime
	content := doc.Find(".t_msgfont").First()
	text := content.Contents().Not(".t_attach").Text()
	html := content.Not(".t_attach").Not("span")
	post.Content = text
	return FixPath(auth), FixPath(title), text, postTime, html, post
}

func FixPath(path string) string {
	path = stringsi.ReplaceRuneEmpty(path, []rune{'\\', '/', ':', ' '})
	if strings.HasSuffix(path, ".") {
		path += "$"
	}
	return path
}

func GetDownloadReq(url, dir string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.NewTaskMeta("Download: " + url),
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return nil, Download(url, dir)
		},
	}
}

func Download(url, dir string) error {

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
	err := client.DownloadImage(dir+name, url)
	if err != nil {
		log.Println(err)
		return DownloadError.Message(url).AppendErr(err)
	}
	return nil
}

func R(url string) (io.Reader, error) {
	var res client.RawResponse
	err := client.NewGetRequest(url).RetryTimes(20).RetryHandle(func(req *client.RequestParams) {
		n := rand.Intn(5)
		req.AddHeader("User-Agent", userAgent[n])
	}).AddHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8").
		AddHeader("Accept-Encoding", "gzip, deflate").
		AddHeader("Accept-Language", "zh-CN,zh;q=0.9;charset=utf-8").
		AddHeader("Connection", "keep-alive").Do(nil, &res)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(res), nil
}

func FixWeb(path string, handle func(int) error) {
	f, err := os.Open(Conf.Pro.CommonDir + path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		id, _ := strconv.Atoi(scanner.Text())
		go handle(id)
		time.Sleep(Conf.Pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

}
