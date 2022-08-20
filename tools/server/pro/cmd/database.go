package main

import (
	"bufio"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	pro.Start(history)
}

func history(sd *pro.Speed) {
	start := 407889
	end := 407912
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetchHistory(i, sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
}

func historyOne(sd *pro.Speed) {
	historyFormFile(`fail_post_2020_12_08_13_07_15`, sd)
}

func historyFormFile(path string, sd *pro.Speed) {
	f, err := os.Open(pro.Conf.Pro.CommonDir + path + pro.Conf.Pro.Ext)
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
			err := pro.Dao.DB.Save(invalidPost).Error
			if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
				sd.FailDB <- tid + " 2"
			}
			continue
		}
		sd.WebAdd(1)
		go fetchHistory(id, sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func fetchHistory(id int, sd *pro.Speed) {
	defer sd.WebDone()
	tid := strconv.Itoa(id)
	reader, err := pro.R(pro.Conf.Pro.CommonUrl + tid)
	if err != nil {
		//log.Println(err, "id:", tid)
		if !strings.HasPrefix(err.Error(), "返回错误") {
			sd.Fail <- tid
		}
		invalidPost := &pro.Post{TId: id, Status: 2}
		err := pro.Dao.DB.Save(invalidPost).Error
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

	s := doc.Find(`img[src="images/common/none.gif"]`)
	_, _, _, _, _, post := pro.ParseHtml(doc)
	post.TId = id
	post.PicNum = uint32(s.Length())
	status := "0"
	if post.PicNum == 0 {
		status = "1"
	}
	err = pro.Dao.DB.Save(post).Error
	if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
		sd.FailDB <- tid + " " + status
	}
}
