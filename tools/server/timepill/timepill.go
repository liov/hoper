package timepill

import (
	"errors"
	postgres "github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	surl "net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"tools/timepill/model"
	"tools/timepill/rpc"
)

func RecordTask() {
	todayDiaries, err := ApiService.GetTodayDiaries(1, 20, "")
	if err != nil {
		log.Error(err)
		return
	}
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
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM `+model.NoteBookTableName+` WHERE id = ?  LIMIT 1)`, notebookId).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	if !exists {
		notebook, err := ApiService.GetNotebook(notebookId)
		if err != nil {
			log.Error(err)
			return
		}
		if notebook != nil && notebook.Id > 0 {
			Dao.Hoper.Create(notebook)
		}
	}

}

func RecordDiary(diary *model.Diary) {
	if diary == nil || diary.Id == 0 || DiaryExists(diary.Id) {
		return
	}

	err := Dao.Hoper.Create(diary).Error
	if err != nil {
		log.Error(err)
	}
	if diary.User != nil && diary.UserId > 0 {
		RecordUser(diary.UserId, diary.User.Name)
	}

	if diary.CommentCount > 0 {
		RecordComment(diary.Id)
	}

	if diary.PhotoUrl != "" {
		err = DownloadPic(diary.UserId, diary.PhotoUrl, diary.Created)
		//err = tnsq.PublishPic(Dao.NsqP.Producer, diary.UserId, diary.PhotoUrl, diary.Created)
		if err != nil {
			log.Error(err)
		}
		go func() {
			rep := rpc.FaceDetection(diary.PhotoUrl)
			if rep.Found {
				err = Dao.Hoper.Create(&model.Face{
					UserId:  diary.UserId,
					DairyId: diary.Id,
				}).Error
				if err != nil {
					log.Error(err)
				}
			}
		}()
	}
}

func RecordDiaryById(diaryId int) {
	diary, err := ApiService.GetDiary(diaryId)
	if err != nil {
		log.Error(err)
		return
	}
	RecordDiary(diary)
}

func TodayRecord() {
	var page = 1
	for {
		todayDiaries, err := ApiService.GetTodayDiaries(page, 20, "")
		if err != nil {
			log.Error(err)
			break
		}
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

func DownloadPic(userId int, url, created string) error {
	if url == "" {
		return errors.New("url is empty")
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
		return err
	}
	return CopyDatePic(filepath, date, strconv.Itoa(userId), path.Base(suffixpath))

}

func CopyDatePic(filepath, date, userId, filename string) error {
	dir := Conf.TimePill.PhotoPath + "/"
	year := date[0:4] + "_"
	_, err := os.Stat(dir + year + "/" + date)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir+year+"/"+date, 0666)
		if err != nil {
			return err
		}
	}
	return fs.CopyFile(filepath, dir+year+"/"+date+"/"+userId+"_"+filename)
}

func DownloadCover(typ, url string) error {
	if url == "" || strings.HasSuffix(url, "default.jpg") {
		return errors.New("url is empty")
	}
	URL, _ := surl.Parse(url)
	v := URL.Query().Get("v")
	filepath := Conf.TimePill.PhotoPath
	filename := strings.Split(URL.Path, typ)
	var originFileName string
	if len(filename) > 1 {
		originFileName = filename[1]
	} else {
		originFileName = path.Base(URL.Path)
	}
	originFileName = strings.TrimSuffix(originFileName, path.Ext(originFileName)) + "-v" + v + path.Ext(originFileName)
	filepath += typ + originFileName
	return client.DownloadImage(filepath, url)
}

func DiaryExists(diaryId int) bool {
	exists, err := postgres.Exists(Dao.Hoper.DB, model.DiaryTableName, "id ", diaryId, false)
	if err != nil {
		log.Error(err)
	}
	return exists
}

func UserExists(userId int) bool {
	exists, err := postgres.Exists(Dao.Hoper.DB, model.UserTableName, "user_id ", userId, false)
	if err != nil {
		log.Error(err)
	}
	return exists
}

func UserExistsByIdName(userId int, userName string) bool {
	var exists bool
	err := Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM `+model.UserTableName+` WHERE user_id = ? AND name = ? LIMIT 1)`, userId, userName).Row().Scan(&exists)
	if err != nil {
		log.Error(err)
	}
	return exists
}

func RecordUserDiaries(user *model.User) {
	notebooks, err := ApiService.GetUserNotebooks(user.UserId)
	if err != nil {
		log.Error(err)
		return
	}
	for _, nodebook := range notebooks {
		Dao.Hoper.Create(nodebook)
		if nodebook.CoverUrl != "" {
			err = DownloadCover(model.BookCoverType.String(), nodebook.CoverUrl)
			//err = tnsq.PublishCover(Dao.NsqP.Producer, model.BookCoverType, nodebook.CoverUrl)
			if err != nil {
				log.Error(err)
			}
		}
		var page = 1
		for {
			diaries, err := ApiService.GetNotebookDiaries(nodebook.Id, page, 20)
			if err != nil {
				log.Error(err)
				break
			}
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
	comments, err := ApiService.GetDiaryComments(diaryId)
	if err != nil {
		log.Error(err)
	}
	for _, comment := range comments {
		RecordUser(comment.UserId, comment.User.Name)
		Dao.Hoper.Create(comment)
	}
}

func RecordCommentWithJudge(diaryId int) {
	comments, err := ApiService.GetDiaryComments(diaryId)
	if err != nil {
		log.Error(err)
	}
	for _, comment := range comments {
		if exists, _ := postgres.ExistsById(Dao.Hoper.DB, model.CommentTableName, uint64(comment.Id)); exists {
			continue
		}
		RecordUser(comment.UserId, comment.User.Name)
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

func RecordUserById(userId int) *model.User {
	user, err := ApiService.GetUserInfo(userId)
	if err != nil {
		log.Error(err)
	}
	if user != nil && user.UserId > 0 {
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
			err = DownloadCover(model.UserCoverType.String(), user.CoverUrl)
			//err = tnsq.PublishCover(Dao.NsqP.Producer, model.UserCoverType, user.CoverUrl)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return user
}

func TodayCommentRecord() {
	var page = 1
	today := time.Now().Format("2006-01-02")
	for {
		var diaryIds []int
		err := Dao.Hoper.Table(model.DiaryTableName).Where(`created > ?`, today).Order(`id`).Offset((page-1)*100).Limit(100).Pluck("id", &diaryIds)
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

func CronCommentRecord() {
	var page = 1
	today := time.Now().Format("2006-01-02")
	for {
		var diaryIds []int
		err := Dao.Hoper.Table(model.DiaryTableName).Where(`created > ?`, today).Order(`id`).Offset((page-1)*100).Limit(100).Pluck("id", &diaryIds)
		if err != nil {
			log.Error(err)
		}
		for _, id := range diaryIds {
			RecordCommentWithJudge(id)
		}
		if len(diaryIds) < 100 {
			return
		}
		page++
	}
}

func RecordByNoteBookId(id int) *model.NoteBook {
	page, pageNum := 1, 20
	notebook, err := ApiService.GetNotebook(id)
	if err != nil {
		log.Error(err)
	}
	if notebook.Id == 0 {
		return notebook
	}
	Dao.Hoper.Create(&notebook)
	user, err := ApiService.GetUserInfo(notebook.UserId)
	if err != nil {
		log.Error(err)
	}
	for {
		diaries, err := ApiService.GetNotebookDiaries(id, page, pageNum)
		if err != nil {
			log.Error(err)
			break
		}
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
