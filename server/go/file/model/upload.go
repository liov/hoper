package model

import (
	"errors"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
	"time"
)

type Upload struct {
	Id        int       `gorm:"primary_key" json:"id"`
	FileId    string    `gorm:"index" json:"file_id"`
	File      File      `gorm:"foreignkey:FileId" json:"file"`
	UserId    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    uint8     `gorm:"type:smallint;default:0" json:"status"`
}

type File struct {
	ID   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	MD5  string `gorm:"type:varchar(32)" json:"md5"`
	Ext  string `json:"ext"`
	Path string `json:"path"`
	// Total file size in bytes specified in the NewUpload call
	Size int64 `json:"size"`
	// Indicates whether the total file size is deferred until later
	SizeIsDeferred bool `json:"sizeIsDeferred"`
	// Offset in bytes (zero-based)
	Offset   int64             `json:"offset"`
	MetaData map[string]string `gorm:"type:jsonb" json:"metadata"`
	// Indicates that this is a partial upload which will later be used to form
	// a final upload by concatenation. Partial uploads should not be processed
	// when they are finished since they are only incomplete chunks of files.
	IsPartial bool `json:"isPartial"`
	// Indicates that this is a final upload
	IsFinal bool `json:"isFinal"`
	// If the upload is a final one (see IsFinal) this will be a non-empty
	// ordered slice containing the ids of the uploads of which the final upload
	// will consist after concatenation.
	PartialUploads []string          `gorm:"type:text[]" json:"partial_uploads"`
	Storage        map[string]string `gorm:"type:jsonb" json:"storage"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `json:"-"`
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
	Id  string `json:"id"`
	URL string `json:"url"`
}

type MultiRep struct {
	Id      string `json:"id"`
	URL     string `json:"url"`
	Success bool   `json:"success"`
}
