package data

import (
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/liov/hoper/server/go/protobuf/file"
)

func (d uploadDao) FileInfo(md5, size string) (*model.FileInfo, error) {
	var file model.FileInfo
	raw := `SELECT * FROM ` + model.TableNameUploadInfo + ` WHERE md5 = ? AND size = ? LIMIT 1`
	err := d.Raw(raw, md5, size).Scan(&file).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	if file.Id == "" {
		return nil, nil
	}
	return &file, nil
}

func (d uploadDao) GetUrls(ids []uint64) ([]*file.File, error) {
	var uploadInfos []*file.File
	err := d.Table(model.TableNameFileInfo).Where(`id IN (?)`, ids).Find(&uploadInfos).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return uploadInfos, nil
}
func (d uploadDao) GetUrlsByStrId(ids string) ([]*file.File, error) {
	var uploadInfos []*file.File
	err := d.Table(model.TableNameFileInfo).Where(`id IN (` + ids + `)`).Find(&uploadInfos).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return uploadInfos, nil
}
