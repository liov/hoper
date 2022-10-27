package backup

import (
	initpostgres "github.com/actliboy/hoper/server/go/lib/initialize/db/postgres"
	"time"
)

type dao3 struct {
	Hoper initpostgres.DB
}

func (d dao3) Init() {
}

var Dao dao3

type File struct {
	Id         int
	Name       string `json:"name" gorm:"not null;index"`
	Level      int
	Size       int
	Pid        int
	CreateTime time.Time `gorm:"not null"`
	ModTime    time.Time `gorm:"not null;index"`
}

func (receiver *File) TableName() string {
	return "bak_dir1"
}
