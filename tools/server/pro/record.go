package pro

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	py2 "github.com/actliboy/hoper/server/go/lib/utils/strings/pinyin"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Record(id int, sd *Speed) string {
	tid := strconv.Itoa(id)
	reader, err := Request(http.DefaultClient, Conf.Pro.CommonUrl+tid)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "返回错误") {
			log.Println(err)
		}
		return ""
	}
	defer reader.Close()
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
		return ""
	}

	s := doc.Find(`img[src="images/common/none.gif"]`)

	auth, title, text, postTime, htl, post := ParseHtml(doc)
	post.TId = id
	post.PicNum = uint32(s.Length())
	status := "0"
	if post.PicNum == 0 {
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

	post.Path = dir[CommonDirLen-7:]
	err = Dao.DB.Save(post).Error
	if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
		sd.FailDB <- tid + " " + status
	}

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			log.Println(err, dir)
			sd.Fail <- tid
			return ""
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

	s.Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("file"); ok {
			sd.Add(1)
			go Download(url, dir, sd)
			time.Sleep(Conf.Pro.Interval)
		}
	})
	return tid
}

func (f Fail) Record(name string) {
	go func() {
		file, _ := os.Create(Conf.Pro.CommonDir + name + time.Now().Format("2006_01_02_15_04_05") + Conf.Pro.Ext)
		for txt := range f {
			file.WriteString(txt + "\n")
		}
		file.Close()
	}()
}
