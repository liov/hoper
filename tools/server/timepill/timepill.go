package timepill

import (
	"encoding/base64"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache_ristretto"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/elastic"
	initializeredis "github.com/actliboy/hoper/server/go/lib/tiga/initialize/redis"
	gormi "github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"

	"io"
	surl "net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"tools/toshi"
)

type Config struct {
	TimePill Customize
	Search   toshi.Config
}

func (c Config) Init() {
	Token = "Basic " + base64.StdEncoding.EncodeToString([]byte(Conf.TimePill.User+":"+Conf.TimePill.Password))
}

type Customize struct {
	User        string
	Password    string
	PhotoPath   string
	PhotoPrefix string
	SearchHost  string
	PageSize    int
	Timer       time.Duration
}

type TimepillDao struct {
	Hoper     db.DB
	Redis     initializeredis.Redis
	Cache     cache_ristretto.Cache
	UserCache cache_ristretto.Cache
	Es        elastic.Es `config:"elasticsearch"`
}

func (dao TimepillDao) Init() {
}

func (dao TimepillDao) Close() {
}

var (
	Dao   TimepillDao
	Conf  Config
	Token string
)

func StartRecord() {
	RecordTask()
	tc := time.NewTicker(time.Second * Conf.TimePill.Timer)
	for {
		select {
		case <-tc.C:
			RecordTask()
		}
	}
}

func RecordTask() {
	todayDiaries := ApiService.GetTodayDiaries(1, 20, "")
	for _, diary := range todayDiaries.Diaries {
		if _, ok := Dao.Cache.Get(diary.Id); ok {
			continue
		}
		Dao.Cache.Set(diary.Id, diary.Id, 1)
		Dao.Cache.Wait()
		RecordNoteBook(diary.NoteBookId)
		RecordDiary(diary)
	}
}

func RecordNoteBook(notebookId int) {
	var exists bool
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "note_book" WHERE id = ?  LIMIT 1)`, notebookId).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	if !exists {
		notebook := ApiService.GetNotebook(notebookId)
		if notebook.Id != 0 {
			Dao.Hoper.Create(notebook)
		}
	}

}

func RecordDiary(diary *Diary) {
	if diary == nil || DiaryExists(diary.Id) {
		return
	}

	err := Dao.Hoper.Create(diary).Error
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

func RecordDiaryById(diaryId int) {
	diary := ApiService.GetDiary(diaryId)
	if diary != nil {
		RecordDiary(diary)
	}
}

func TodayRecord() {
	var page = 1
	for {
		todayDiaries := ApiService.GetTodayDiaries(page, 20, "")
		for _, diary := range todayDiaries.Diaries {
			RecordNoteBook(diary.NoteBookId)
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
	if strings.HasSuffix(suffixpath, ".") {
		suffixpath += "jpg"
	}
	if strings.Contains(suffixpath, "image") {
		substr := strings.Split(suffixpath, ".")
		suffixpath = substr[0] + ".jpg"
	}
	date := suffixpath[1:11]
	if strings.Contains(date, "/") {
		date = created[0:10]
	}
	num := userId / 10000
	prepath := Conf.TimePill.PhotoPath + "/" + strconv.Itoa(num) + "-" + strconv.Itoa(num+1) + "/" + strconv.Itoa(userId) + "/"
	filepath := prepath + date + "_" + path.Base(suffixpath)

	err := client.DownloadImage(filepath, url)
	if err != nil {
		log.Error(err)
	} else {
		CopyDatePic(filepath, date, strconv.Itoa(userId), path.Base(suffixpath))
	}
}

func CopyDatePic(filepath, date, userId, filename string) {
	dir := Conf.TimePill.PhotoPath + "/"
	year := date[0:4] + "_"
	_, err := os.Stat(dir + year + "/" + date)
	if os.IsNotExist(err) {
		os.MkdirAll(dir+year+"/"+date, 0666)
	}
	dst, err := os.Create(dir + year + "/" + date + "/" + userId + "_" + filename)
	if err != nil {
		log.Error(err)
	}
	src, err := os.Open(filepath)
	if err != nil {
		log.Error(err)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Error(err)
	}
	dst.Close()
	src.Close()
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
		Dao.Hoper.Where(`is_record = ?`, false).First(&user)
		RecordUserDiaries(&user)
	}
}

func DiaryExists(diaryId int) bool {
	var exists bool
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM diary WHERE id = ? LIMIT 1)`, diaryId).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	return exists
}

func UserExists(userId int) bool {
	var exists bool
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "user" WHERE user_id = ? LIMIT 1)`, userId).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	return exists
}

func UserExistsByIdName(userId int, userName string) bool {
	var exists bool
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "user" WHERE user_id = ? AND name = ? LIMIT 1)`, userId, userName).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	return exists
}

func RecordUserDiaries(user *User) {
	notebooks := ApiService.GetUserNotebooks(user.UserId)
	for _, nodebook := range notebooks {
		Dao.Hoper.Create(nodebook)
		if nodebook.CoverUrl != "" {
			DownloadCover(nodebook.CoverUrl)
		}
		var page = 1
		for {
			diaries := ApiService.GetNotebookDiaries(nodebook.Id, page, 20)
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
}

func RecordComment(diaryId int) {
	comments := ApiService.GetDiaryComments(diaryId)
	for _, comment := range comments {
		RecordUser(comment.UserId, comment.User.Name)
		Dao.Hoper.Create(comment)
	}
}

func RecordCommentWithJudge(diaryId int) {
	comments := ApiService.GetDiaryComments(diaryId)
	for _, comment := range comments {
		RecordUser(comment.UserId, comment.User.Name)
		if exists, _ := gormi.ExistsById(Dao.Hoper.DB, "comment", uint64(comment.Id)); exists {
			continue
		}
		Dao.Hoper.Create(comment)
	}
}

func RecordUser(userId int, userName string) {
	var exists bool
	if _, ok := Dao.UserCache.Get(userId); ok {
		exists = true
	}
	Dao.UserCache.Set(userId, userId, 1)
	Dao.UserCache.Wait()
	if !exists && !UserExistsByIdName(userId, userName) {
		RecordUserById(userId)
	}
}

func RecordUserById(userId int) *User {
	user := ApiService.GetUserInfo(userId)
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
	return user
}

func TodayCommentRecord() {
	var page = 1
	today := time.Now().Format("2006-01-02")
	for {
		var diaryIds []int
		err := Dao.Hoper.Table(`diary`).Where(`created > ?`, today).Order(`id`).Offset((page-1)*100).Limit(100).Pluck("id", &diaryIds)
		if err != nil {
			log.Error(err)
		}
		for _, id := range diaryIds {
			RecordComment(id)
		}
		if len(diaryIds) < 100 {
			return
		}
		page++
	}
}

func RecordByNoteBookId(id int) *NoteBook {
	page, pageNum := 1, 20
	notebook := ApiService.GetNotebook(id)
	if notebook.Id == 0 {
		return notebook
	}
	Dao.Hoper.Create(&notebook)
	user := ApiService.GetUserInfo(notebook.UserId)
	for {
		diaries := ApiService.GetNotebookDiaries(id, page, pageNum)
		if diaries.Items == nil {
			break
		}
		for _, diary := range diaries.Items {
			diary.User = user
			RecordDiary(diary)
		}
		if len(diaries.Items) < pageNum {
			break
		}
		page++
	}
	return notebook
}
