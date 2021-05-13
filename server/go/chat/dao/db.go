package dao

import (
	"github.com/liov/hoper/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/v2/upload/model"
	"gorm.io/gorm"
)

func (d *uploadDao) UploadDB(db *gorm.DB, md5, size string) (*model.UploadInfo, error) {
	var upload model.UploadInfo
	raw := `SELECT * FROM ` + model.UploadTableName + ` WHERE md5 = ? AND size = ? LIMIT 1`
	err := db.Raw(raw, md5, size).Scan(&upload).Error
	if err != nil {
		return nil, d.ctxi.ErrorLog(errorcode.DBError, err, "FileExistsDB")
	}
	if upload.Id == 0 {
		return nil, nil
	}
	return &upload, nil
}
