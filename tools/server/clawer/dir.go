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

type Dir struct {
	Type      int       `json:"type"`
	Date      string    `json:"date" gorm:"type:date"`
	UserId    int       `json:"userId"`
	KeyId     int       `json:"keyId"`
	KeyIdStr  string    `json:"keyIdStr"`
	BaseUrl   string    `json:"baseUrl"`
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
	filepath := strings.Join([]string{d.DownloadPath, d.Date[:4], d.Date[:7], d.Date, strconv.Itoa(d.UserId) + "_" + d.KeyIdStr + "_" + d.BaseUrl}, "/")

	if fs.NotExist(filepath) {
		err := client.DownloadFileWithRefer(filepath, d.Url, d.Referer)
		if err != nil {
			log.Info("下载图片失败：", err)
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
