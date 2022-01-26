package main

import (
	"encoding/base64"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"github.com/dgraph-io/ristretto"
	"gorm.io/gorm"
	"path"
	"strings"
	"time"
)

type Config struct {
	Hoper    db.DatabaseConfig
	Cache    cache.CacheConfig
	TimePill Customize
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
	Hoper *gorm.DB
	Cache *ristretto.Cache
}

func (d Dao) Init() {
}

func (d Dao) Close() {
}

var (
	dao   Dao
	conf  Config
	token string
)

func main() {
	defer initialize.Start(&conf, &dao)()
	token = "Basic " + base64.StdEncoding.EncodeToString([]byte(conf.TimePill.User+":"+conf.TimePill.Password))
	todayRecord()
	record()
	tc := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-tc.C:
			record()
		}
	}
}

func createTable() {
	fmt.Println(dao.Hoper.Migrator().CreateTable(&Diary{}))
}

func record() {
	todayDiaries := getTodayDiaries(1, 20, "")
	for _, diary := range todayDiaries.Diaries {
		if _, ok := dao.Cache.Get(diary.Id); ok {
			continue
		}
		dao.Cache.Set(diary.Id, diary.Id, 1)
		dao.Cache.Wait()
		err := dao.Hoper.Create(diary).Error
		if err != nil {
			log.Error(err)
		}
		if diary.PhotoUrl != "" {
			err = downloadPic(diary.PhotoUrl)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func todayRecord() {
	var page = 1
	for {
		todayDiaries := getTodayDiaries(page, 20, "")
		for _, diary := range todayDiaries.Diaries {
			err := dao.Hoper.Create(diary).Error
			if err != nil {
				log.Error(err)
			}
			if diary.PhotoUrl != "" {
				err = downloadPic(diary.PhotoUrl)
				if err != nil {
					log.Error(err)
				}
			}
		}
		if len(todayDiaries.Diaries) < 20 {
			return
		}
		page++
	}
}

func downloadPic(url string) error {
	if url == "" {
		return nil
	}
	filepath := conf.TimePill.PhotoPath
	filename := strings.Split(url, "photos")
	if len(filename) > 1 {
		filepath += filename[1]
	} else {
		filepath += "/" + time.Now().Format("2006-01-02") + path.Base(url)
	}

	return client.DownloadImage(filepath, url)
}
