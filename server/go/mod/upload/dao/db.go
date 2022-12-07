package dao

import (
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	"github.com/liov/hoper/server/go/mod/protobuf/upload"
	"github.com/liov/hoper/server/go/mod/upload/model"
	"gorm.io/gorm"
)

func (d *uploadDao) UploadDB(db *gorm.DB, md5, size string) (*model.UploadInfo, error) {
	var uploadInfo model.UploadInfo
	raw := `SELECT * FROM ` + model.UploadTableName + ` WHERE md5 = ? AND size = ? LIMIT 1`
	err := db.Raw(raw, md5, size).Scan(&uploadInfo).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError, err, "FileExistsDB")
	}
	if uploadInfo.Id == 0 {
		return nil, nil
	}
	return &uploadInfo, nil
}

func (d *uploadDao) GetUrls(db *gorm.DB, ids []uint64) ([]*upload.UploadInfo, error) {
	var uploadInfos []*upload.UploadInfo
	err := db.Table(model.UploadTableName).Where(`id IN (?)`, ids).Find(&uploadInfos).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError, err, "GetUrls")
	}
	return uploadInfos, nil
}
func (d *uploadDao) GetUrlsByStrId(db *gorm.DB, ids string) ([]*upload.UploadInfo, error) {
	var uploadInfos []*upload.UploadInfo
	err := db.Table(model.UploadTableName).Where(`id IN (` + ids + `)`).Find(&uploadInfos).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError, err, "GetUrlsByStrId")
	}
	return uploadInfos, nil
}
