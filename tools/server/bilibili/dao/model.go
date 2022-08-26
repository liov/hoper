package dao

import (
	"time"
)

const (
	Schema           = "bilibili."
	TableNameView    = Schema + "view"
	TableNameViewBak = Schema + "view_bak"
	TableNameVideo   = Schema + "video"
)

type View struct {
	Bvid        string `json:"bvid" gorm:"index:idx_bvid,unique"`
	Aid         int    `json:"aid" gorm:"primaryKey"`
	Data        []byte `json:"data" gorm:"type:jsonb"`
	CoverRecord bool
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
}

func (v *View) TableName() string {
	return TableNameView
}

type ViewBak struct {
	Id        int       `gorm:"primaryKey"`
	Aid       int       `json:"aid" gorm:"index"`
	Data      []byte    `json:"data" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
}

func (v *ViewBak) TableName() string {
	return TableNameViewBak
}

type Video struct {
	Aid       int    `json:"aid" gorm:"primaryKey"`
	Cid       int    `json:"cid" gorm:"primaryKey"`
	Data      []byte `json:"data" gorm:"type:jsonb"`
	Record    bool
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (v *Video) TableName() string {
	return TableNameVideo
}
