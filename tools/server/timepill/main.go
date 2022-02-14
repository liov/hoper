package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"github.com/dgraph-io/ristretto"
	"gorm.io/gorm"
	surl "net/url"
	"path"
	"strings"
	"time"
)

type Config struct {
	Hoper     db.DatabaseConfig
	Cache     cache.CacheConfig
	UserCache cache.CacheConfig
	TimePill  Customize
}

func (c Config) Init() {
}

type Customize struct {
	User        string
	Password    string
	PhotoPath   string
	PhotoPrefix string
}

type Dao struct {
	Hoper     *gorm.DB
	Cache     *ristretto.Cache
	UserCache *ristretto.Cache
}

func (d Dao) Init() {
}

func (d Dao) Close() {
}

var (
	dao     Dao
	conf    Config
	token   string
	picChan = make(chan struct{}, 10)
	today   = flag.Bool("t", false, "记录今天日记")
)

func main() {
	defer initialize.Start(&conf, &dao)()
	token = "Basic " + base64.StdEncoding.EncodeToString([]byte(conf.TimePill.User+":"+conf.TimePill.Password))
	go recordByUser()
	if *today {
		fmt.Println("todayRecord")
		todayRecord()
	}
	fmt.Println("startRecord")
	startRecord()
}

func createTable() {
	fmt.Println(dao.Hoper.Migrator().CreateTable(&Badge{}))
}

func startRecord() {
	record()
	tc := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-tc.C:
			record()
		}
	}
}

func record() {
	todayDiaries := getTodayDiaries(1, 20, "")
	for _, diary := range todayDiaries.Diaries {
		if _, ok := dao.Cache.Get(diary.Id); ok {
			continue
		}
		dao.Cache.Set(diary.Id, diary.Id, 1)
		dao.Cache.Wait()
		recordDiary(diary)
	}
}

func recordDiary(diary *Diary) {
	if diary == nil {
		return
	}
	var exists bool
	err := dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM diary WHERE id = ? LIMIT 1)`, diary.Id).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	if exists {
		return
	}
	err = dao.Hoper.Create(diary).Error
	if err != nil {
		log.Error(err)
	}

	recordUser(diary.UserId, diary.User.Name)

	if diary.CommentCount > 0 {
		recordComment(diary.Id)
	}

	if diary.PhotoUrl != "" {
		picChan <- struct{}{}
		go downloadPic(diary.PhotoUrl)
	}
}

func todayRecord() {
	var page = 1
	for {
		todayDiaries := getTodayDiaries(page, 20, "")
		for _, diary := range todayDiaries.Diaries {
			recordDiary(diary)
		}
		if len(todayDiaries.Diaries) < 20 {
			return
		}
		page++
	}
}

func downloadPic(url string) {
	if url == "" {
		return
	}
	URL, _ := surl.Parse(url)
	var largepath string
	if strings.HasSuffix(URL.Path, "!large") {
		largepath = strings.TrimSuffix(URL.Path, "!large")
	}
	filepath := conf.TimePill.PhotoPath
	filename := strings.Split(URL.Path, "photos")
	if len(filename) > 1 {
		filepath += filename[1]
	} else if largepath != "" {
		filepath += largepath
	} else {
		filepath += "/" + time.Now().Format("2006-01-02") + "/" + path.Base(url)
	}

	err := client.DownloadImage(filepath, url)
	if err != nil {
		log.Error(err)
	}
	<-picChan
}

func downloadCover(url string) {
	if url == "" || strings.HasSuffix(url, "default.jpg") {
		return
	}
	URL, _ := surl.Parse(url)
	filepath := conf.TimePill.PhotoPath
	filename := strings.Split(URL.Path, "book_cover")
	if len(filename) > 1 {
		filepath += "/book_cover" + filename[1]
	} else {
		filepath += "/book_cover/" + path.Base(URL.Path)
	}

	err := client.DownloadImage(filepath, url)
	if err != nil {
		log.Error(err)
	}
	<-picChan
}

func downloadUserCover(url string) {
	if url == "" || strings.HasSuffix(url, "default.jpg") {
		return
	}
	URL, _ := surl.Parse(url)
	filepath := conf.TimePill.PhotoPath
	filename := strings.Split(URL.Path, "user_icon")
	if len(filename) > 1 {
		filepath += "/user_icon" + filename[1]
	} else {
		filepath += "/user_icon/" + path.Base(URL.Path)
	}

	err := client.DownloadImage(filepath, url)
	if err != nil {
		log.Error(err)
	}
	<-picChan
}

func recordByUser() {
	for {
		var user User
		err := dao.Hoper.Where(`is_record = ?`, false).First(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		}
		recordUserDiaries(&user)
	}
}

func recordUserDiaries(user *User) {
	notebooks := getUserNotebooks(user.UserId)
	for _, nodebook := range notebooks {
		dao.Hoper.Create(nodebook)
		if nodebook.CoverUrl != "" {
			picChan <- struct{}{}
			go downloadCover(nodebook.CoverUrl)
		}
		var page = 1
		for {
			diaries := getNotebookDiaries(nodebook.Id, page, 20)
			for _, diary := range diaries.Items {
				diary.User = user
				recordDiary(diary)
			}
			if len(diaries.Items) < 20 {
				break
			}
			page++
		}
	}
	dao.Hoper.Table("user").Where("user_id = ?", user.UserId).Update("is_record", true)
}

func recordComment(diaryId int) {
	comments := getDiaryComments(diaryId)
	for _, comment := range comments {
		recordUser(comment.UserId, comment.User.Name)
		dao.Hoper.Create(comment)
	}
}

func fixNoteBook() {
	var users []*User
	dao.Hoper.Where(`is_record = ?`, true).Find(&users)
	for _, user := range users {
		notebooks := getUserNotebooks(user.UserId)
		for _, nodebook := range notebooks {
			dao.Hoper.Create(nodebook)
		}
	}

}

func recordUser(userId int, userName string) {
	var exists bool
	if _, ok := dao.UserCache.Get(userId); ok {
		exists = true
	}
	dao.UserCache.Set(userId, userId, 1)
	dao.UserCache.Wait()
	if !exists {
		err := dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "user" WHERE user_id = ? AND name = ? LIMIT 1)`, userId, userName).Row().Scan(&exists)
		if err != nil {
			log.Error(err)
		}
	}
	if !exists {
		user := getUserInfo(userId)
		if user.UserId != 0 {
			err := dao.Hoper.Create(user).Error
			if err != nil {
				log.Error(err)
			}
			if user.Badges != nil {
				for _, bage := range user.Badges {
					dao.Hoper.Create(bage)
				}
			}
			if user.CoverUrl != "" {
				picChan <- struct{}{}
				go downloadUserCover(user.CoverUrl)
			}
			if user.IconUrl != "" && user.IconUrl != user.CoverUrl {
				picChan <- struct{}{}
				go downloadUserCover(user.IconUrl)
			}
		}
	}
}
