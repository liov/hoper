package claweri

import (
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Platform int    `json:"platform"`
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	PicNums  int    `json:"picNums"`
	Sort     int    `json:"sort"`
}

type Dir struct {
	Platform  int       `json:"platform"`
	UserId    int       `json:"userId"`
	KeyId     int       `json:"keyId"`
	KeyIdStr  string    `json:"keyIdStr"`
	BaseUrl   string    `json:"baseUrl"`
	Type      int       `json:"type"`
	PubAt     string    `json:"pubAt" gorm:"type:timestamptz(0);default:0001-01-01 00:00:00"`
	CreatedAt time.Time `json:"createdAt"`
}

func (d *Dir) TableName() string {
	return "dir"
}

type DownloadMeta struct {
	Dir
	DownloadPath string `json:"downloadPath"`
	Url          string `json:"url"`
	Referer      string `json:"referer"`
}

func (d *DownloadMeta) Download(db *gorm.DB) error {
	if d.KeyId != 0 {
		d.KeyIdStr = strconv.Itoa(d.KeyId)
	}
	userIdStr := strconv.Itoa(d.UserId)
	filepath := strings.Join([]string{d.DownloadPath, userIdStr, d.PubAt[:4], d.PubAt[:7], d.PubAt[:10], userIdStr + "_" + d.KeyIdStr + "_" + d.PubAt[:10] + "-" + d.PubAt[11:] + "_" + d.BaseUrl}, "/")
	var err error
	if fs.NotExist(filepath) {
		if d.Referer != "" {
			err = client.DownloadFileWithRefer(filepath, d.Url, d.Referer)
		} else {
			err = client.DownloadFile(filepath, d.Url)
		}
		if err != nil {
			log.Info("下载文件失败：", err)
			return err
		}
		err = db.Create(&d.Dir).Error
		if err != nil {
			return err
		}
		log.Info("下载文件成功：", filepath)
	} else {
		log.Info("文件已存在：", filepath)
	}
	return nil
}
