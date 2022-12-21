package claweri

import (
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
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
	Platform int    `json:"platform" comment:"1-pro,2-timepill,3-bilibili,4-weibo"`
	UserId   int    `json:"userId"`
	KeyId    int    `json:"keyId"`
	KeyIdStr string `json:"keyIdStr"`
	//PrePath   string    `json:"prePath" comment:""`
	BaseUrl   string    `json:"baseUrl"`
	Type      int       `json:"type" comment:"1-图片，2-live图片，3-视频，4-动图"`
	PubAt     time.Time `json:"pubAt" gorm:"type:timestamptz(0);default:0001-01-01 00:00:00"`
	CreatedAt time.Time `json:"createdAt"`
}

func (d *Dir) TableName() string {
	return "dir"
}

func (d *Dir) Path() string {
	if d.KeyId != 0 && d.KeyIdStr == "" {
		d.KeyIdStr = strconv.Itoa(d.KeyId)
	}
	date := d.PubAt.Format(timei.TimeFormatCompact)
	compactPubAt := stringsi.ReplaceRuneEmpty(date, '-', ' ', ':')
	userIdStr := strconv.Itoa(d.UserId)
	filepath := strings.Join([]string{userIdStr, date[:4], compactPubAt + "_" + userIdStr + "_" + d.KeyIdStr + "_" + d.BaseUrl}, "/")
	return filepath
}

type DownloadMeta struct {
	Dir
	DownloadPath string `json:"downloadPath"`
	Url          string `json:"url"`
	Referer      string `json:"referer"`
}

func (d *DownloadMeta) Download(db *gorm.DB) error {
	filepath := d.DownloadPath + "/" + d.Path()
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
