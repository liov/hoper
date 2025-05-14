package data

import (
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/liov/hoper/server/go/protobuf/upload"
	"gorm.io/gorm"
)

func (d *uploadDao) FileInfo(db *gorm.DB, md5, size string) (*model.File, error) {
	var file model.File
	raw := `SELECT * FROM ` + model.TableNameUpload + ` WHERE md5 = ? AND size = ? LIMIT 1`
	err := db.Raw(raw, md5, size).Scan(&file).Error
	if err != nil {
		return nil, d.RespErrorLog(errcode.DBError, err, "FileExistsDB")
	}
	if file.Id == "" {
		return nil, nil
	}
	return &file, nil
}

func (d *uploadDao) GetUrls(db *gorm.DB, ids []uint64) ([]*upload.File, error) {
	var uploadInfos []*upload.File
	err := db.Table(model.TableNameFile).Where(`id IN (?)`, ids).Find(&uploadInfos).Error
	if err != nil {
		return nil, d.RespErrorLog(errcode.DBError, err, "GetUrls")
	}
	return uploadInfos, nil
}
func (d *uploadDao) GetUrlsByStrId(db *gorm.DB, ids string) ([]*upload.File, error) {
	var uploadInfos []*upload.File
	err := db.Table(model.TableNameFile).Where(`id IN (` + ids + `)`).Find(&uploadInfos).Error
	if err != nil {
		return nil, d.RespErrorLog(errcode.DBError, err, "GetUrlsByStrId")
	}
	return uploadInfos, nil
}
