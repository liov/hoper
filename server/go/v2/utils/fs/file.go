package fs

import (
	"errors"
	"mime/multipart"
	"strings"
)

type File struct {
	ID           uint64 `gorm:"primary_key" json:"id"`
	FileName     string `gorm:"type:varchar(100);not null" json:"file_name"`
	OriginalName string `gorm:"type:varchar(100);not null" json:"original_name"`
	URL          string `json:"url"`
	MD5          string `gorm:"type:varchar(32)" json:"md5"`
	Mime         string `json:"mime"`
	Size         uint64 `json:"size"`
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