package data

import (
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/upload/model"
	"gorm.io/gorm"
)

func (d *uploadDao) UploadDB(db *gorm.DB, md5, size string) (*model.UploadInfo, error) {
	var upload model.UploadInfo
	raw := `SELECT * FROM ` + model.UploadTableName + ` WHERE md5 = ? AND size = ? LIMIT 1`
	err := db.Raw(raw, md5, size).Scan(&upload).Error
	if err != nil {
		return nil, d.RespErrorLog(errcode.DBError, err, "FileExistsDB")
	}
	if upload.Id == 0 {
		return nil, nil
	}
	return &upload, nil
}
