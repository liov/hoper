package clawer

import timei "github.com/liov/hoper/server/go/lib/utils/time"

type Dir struct {
	Type     int        `json:"type"`
	UserId   int        `json:"userId"`
	UserName string     `json:"userName"`
	Date     timei.Date `json:"date" gorm:"type:"`
	FilePath string     `json:"filePath"`
}
