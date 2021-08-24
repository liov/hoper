package model

import (
	"errors"
	"mime/multipart"
	"strings"
	"time"
)

type UploadInfo struct {
	File
	Path      string    `json:"path"`
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    uint8     `gorm:"type:smallint;default:0" json:"status"`
}

type File struct {
	Id   uint64 `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	MD5  string `gorm:"type:varchar(32)" json:"md5"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

func GetExt(file *multipart.FileHeader) (string, error) {
	var ext string
	var index = strings.LastIndex(file.Filename, ".")
	if index == -1 {
		return "", nil
	} else {
		ext = file.Filename[index:]
	}
	if len(ext) == 1 {
		return "", errors.New("无效的扩展名")
	}
	return ext, nil
}

type Rep struct {
	Id  uint64 `json:"id"`
	URL string `json:"url"`
}

type MultiRep struct {
	Id      int    `json:"id"`
	URL     string `json:"url"`
	Success bool   `json:"success"`
}

type UploadExt struct {
	Id        uint64    `gorm:"primary_key" json:"id"`
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UploadId  uint64    `json:"upload_id"`
}
