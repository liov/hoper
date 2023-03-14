package pro

import (
	errorsi "github.com/liov/hoper/server/go/lib/utils/errors"
	"github.com/liov/hoper/server/go/lib/v2/utils/net/http/client/crawler"
	"log"
	"os"
	"time"
)

const (
	DBError errorsi.ErrCode = iota
	ReqPostError
	CreateFileError
	MkdirError
	DownloadError
)

type Post struct {
	ID        uint32
	TId       int    `gorm:"uniqueIndex"`
	UserId    int    `gorm:"default:0;index"`
	UserName  string `gorm:"size:255;default:''"`
	Title     string `gorm:"size:255;default:''"`
	Content   string `gorm:"type:text"`
	CreatedAt string `gorm:"type:timestamptz(6);default:'0001-01-01 00:00:00'"`
	PicNum    uint32 `gorm:"default:0"`
	Score     uint8  `gorm:"default:0"`
	Status    uint8  `gorm:"default:0"`
	Path      string `gorm:"size:255;default:''"`
}

func (p *Post) TableName() string {
	return "post"
}

func ErrorHandler() func(task *crawler.Request) {
	file, _ := os.Create(Conf.Pro.CommonDir + "fail_" + time.Now().Format("2006_01_02_15_04_05") + Conf.Pro.Ext)
	return func(task *crawler.Request) {
		errs := task.Errs()
		err := errs[len(errs)-1]
		log.Println("任务5次错误：", err.Error())
		switch err.(*errorsi.ErrRep).Code {
		case DBError:
			file.WriteString("db " + err.Error() + "\n")
		case DownloadError:
			file.WriteString("pic " + err.Error() + "\n")
		case ReqPostError:
			file.WriteString("req " + err.Error() + "\n")
		}
	}
}
