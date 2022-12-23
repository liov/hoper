package dao

import (
	"github.com/liov/hoper/server/go/lib/utils/dao/db/postgres"
	"time"
)

const (
	Schema           = "bilibili."
	TableNameView    = Schema + "view"
	TableNameViewBak = TableNameView + "_bak"
	TableNameVideo   = Schema + "video"
	TableNameVideoV2 = Schema + "video_v2"
	TableNameUser    = Schema + "user"
)

type User struct {
	Id   int `json:"id" gorm:"primaryKey"`
	Name string
	Face string
}

func (v *User) TableName() string {
	return TableNameUser
}

type View struct {
	Aid       int `json:"aid" gorm:"primaryKey"`
	Bvid      string
	Uid       int
	Title     string
	Desc      string
	Dynamic   string `json:"dynamic"`
	Tid       int
	Pic       string
	Ctime     time.Time
	Tname     string
	Videos    int
	Pubdate   time.Time
	Record    int
	CreatedAt time.Time `json:"created_at" gorm:"index;default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"<-:false;index;"`
}

func (v *View) TableName() string {
	return TableNameView
}

type Video struct {
	Aid           int `json:"aid" gorm:"primaryKey"`
	Cid           int `json:"cid" gorm:"primaryKey"`
	Part          string
	Page          int
	AcceptFormat  string
	VideoCodecid  int
	Duration      int
	AcceptQuality postgres.IntArray `gorm:"type:int2[]"`
	Record        int
	CreatedAt     time.Time `json:"created_at" gorm:"index;default:now()"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"<-:false;index"`
}

func (v *Video) TableName() string {
	return TableNameVideo
}

type ViewBak struct {
	Aid       int `json:"aid" gorm:"primaryKey"`
	Title     string
	Desc      string
	Dynamic   string `json:"dynamic"`
	Tid       string
	Pic       string
	Ctime     time.Time `gorm:"type:timestamptz(6)"`
	Tname     string
	Videos    int
	Pubdate   time.Time `gorm:"type:timestamptz(6)"`
	CreatedAt time.Time `json:"created_at" gorm:"index;default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"<-:false;index"`
}

func (v *ViewBak) TableName() string {
	return TableNameViewBak
}
