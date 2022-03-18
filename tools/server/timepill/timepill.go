package timepill

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache_ristretto"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	initializeredis "github.com/actliboy/hoper/server/go/lib/tiga/initialize/redis"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"gorm.io/gorm"
	surl "net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	TimePill Customize
}

func (c Config) Init() {
	Token = "Basic " + base64.StdEncoding.EncodeToString([]byte(Conf.TimePill.User+":"+Conf.TimePill.Password))
}

type Customize struct {
	User        string
	Password    string
	PhotoPath   string
	PhotoPrefix string
}

type TimepillDao struct {
	Hoper     db.DB
	Redis     initializeredis.Redis
	Cache     cache_ristretto.Cache
	UserCache cache_ristretto.Cache
}

func (d TimepillDao) Init() {
}

func (d TimepillDao) Close() {
}

var (
	Dao   TimepillDao
	Conf  Config
	Token string
)

func CreateTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&Badge{}))
}

func StartRecord() {
	RecordTask()
	tc := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-tc.C:
			RecordTask()
		}
	}
}

func RecordTask() {
	todayDiaries := GetTodayDiaries(1, 20, "")
	for _, diary := range todayDiaries.Diaries {
		if _, ok := Dao.Cache.Get(diary.Id); ok {
			continue
		}
		Dao.Cache.Set(diary.Id, diary.Id, 1)
		Dao.Cache.Wait()
		RecordDiary(diary)
	}
}

func RecordDiary(diary *Diary) {
	if diary == nil {
		return
	}
	var exists bool
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM diary WHERE id = ? LIMIT 1)`, diary.Id).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	if exists {
		return
	}
	err = Dao.Hoper.Create(diary).Error
	if err != nil {
		log.Error(err)
	}

	RecordUser(diary.UserId, diary.User.Name)

	if diary.CommentCount > 0 {
		RecordComment(diary.Id)
	}

	if diary.PhotoUrl != "" {
		DownloadPic(diary.UserId, diary.PhotoUrl, diary.Created)
	}
}

func TodayRecord() {
	var page = 1
	for {
		todayDiaries := GetTodayDiaries(page, 20, "")
		for _, diary := range todayDiaries.Diaries {
			RecordDiary(diary)
		}
		if len(todayDiaries.Diaries) < 20 {
			return
		}
		page++
	}
}

func DownloadPic(userId int, url, created string) {
	if url == "" {
		return
	}
	URL, _ := surl.Parse(url)
	suffixpath := URL.Path
	if strings.HasSuffix(URL.Path, "!large") {
		suffixpath = strings.TrimSuffix(URL.Path, "!large")
	} else if strings.Contains(URL.Path, "photos") {
		suffixpath = strings.Split(URL.Path, "photos")[1]
	}
	date := suffixpath[1:11]
	if strings.Contains(date, "/") {
		date = created[0:10]
	}
	prepath := Conf.TimePill.PhotoPath + "/" + strconv.Itoa(userId) + "/"
	filepath := prepath + date + "-" + path.Base(suffixpath)

	err := client.DownloadImage(filepath, url)
	if err != nil {
		log.Error(err)
	}
}

func DownloadCover(url string) {
	if url == "" || strings.HasSuffix(url, "default.jpg") {
		return
	}
	URL, _ := surl.Parse(url)
	filepath := Conf.TimePill.PhotoPath
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
}

func DownloadUserCover(url string) {
	if url == "" || strings.HasSuffix(url, "default.jpg") {
		return
	}
	URL, _ := surl.Parse(url)
	filepath := Conf.TimePill.PhotoPath
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
}

func RecordByUser() {
	for {
		var user User
		err := Dao.Hoper.Where(`is_record = ?`, false).First(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		}
		RecordUserDiaries(&user)
	}
}

func RecordByOrderUser() {
	ctx := context.Background()
	key := "RecordByOrderUserID"
	err := Dao.Redis.SetNX(ctx, key, 0, 0).Err()
	if err != nil {
		log.Error(err)
	}
	var id int
	err = Dao.Redis.Get(ctx, key).Scan(&id)
	if err != nil {
		log.Error(err)
	}
	tc := time.NewTicker(time.Second * 2)
	for {
		var exists bool
		err = Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "user" WHERE user_id = ? LIMIT 1)`, id).Row().Scan(&exists)
		if err != nil {
			log.Error(err)
		}
		if exists {
			id++
			err = Dao.Redis.Incr(ctx, key).Err()
			if err != nil {
				log.Error(err)
			}
			continue
		}
		<-tc.C
		user := GetUserInfo(id)
		if user.UserId != 0 {
			err = Dao.Hoper.Create(user).Error
			if err != nil {
				log.Error(err)
			}
			if user.Badges != nil {
				for _, bage := range user.Badges {
					Dao.Hoper.Create(bage)
				}
			}
			if user.CoverUrl != "" {
				DownloadUserCover(user.CoverUrl)
			}
			if user.IconUrl != "" && user.IconUrl != user.CoverUrl {
				DownloadUserCover(user.IconUrl)
			}
			RecordUserDiaries(user)
		} else {
			time.Sleep(time.Second)
		}
		id++
		err = Dao.Redis.Incr(ctx, key).Err()
		if err != nil {
			log.Error(err)
		}
	}
}

func RecordUserDiaries(user *User) {
	notebooks := GetUserNotebooks(user.UserId)
	for _, nodebook := range notebooks {
		Dao.Hoper.Create(nodebook)
		if nodebook.CoverUrl != "" {
			DownloadCover(nodebook.CoverUrl)
		}
		var page = 1
		for {
			diaries := GetNotebookDiaries(nodebook.Id, page, 20)
			for _, diary := range diaries.Items {
				diary.User = user
				RecordDiary(diary)
			}
			if len(diaries.Items) < 20 {
				break
			}
			page++
		}
	}
	Dao.Hoper.Table("user").Where("user_id = ?", user.UserId).Update("is_record", true)
}

func RecordComment(diaryId int) {
	comments := GetDiaryComments(diaryId)
	for _, comment := range comments {
		RecordUser(comment.UserId, comment.User.Name)
		Dao.Hoper.Create(comment)
	}
}

func FixNoteBook() {
	var users []*User
	Dao.Hoper.Where(`is_record = ?`, true).Find(&users)
	for _, user := range users {
		notebooks := GetUserNotebooks(user.UserId)
		for _, nodebook := range notebooks {
			Dao.Hoper.Create(nodebook)
		}
	}

}

func RecordUser(userId int, userName string) {
	var exists bool
	if _, ok := Dao.UserCache.Get(userId); ok {
		exists = true
	}
	Dao.UserCache.Set(userId, userId, 1)
	Dao.UserCache.Wait()
	if !exists {
		err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "user" WHERE user_id = ? AND name = ? LIMIT 1)`, userId, userName).Row().Scan(&exists)
		if err != nil {
			log.Error(err)
		}
	}
	if !exists {
		user := GetUserInfo(userId)
		if user.UserId != 0 {
			err := Dao.Hoper.Create(user).Error
			if err != nil {
				log.Error(err)
			}
			if user.Badges != nil {
				for _, bage := range user.Badges {
					Dao.Hoper.Create(bage)
				}
			}
			if user.CoverUrl != "" {
				DownloadUserCover(user.CoverUrl)
			}
			if user.IconUrl != "" && user.IconUrl != user.CoverUrl {
				DownloadUserCover(user.IconUrl)
			}
		}
	}
}

func RecordByDiaryBook() {

}
