package db

import (
	"github.com/hopeio/protobuf/errcode"
	dbi "github.com/hopeio/utils/dao/database"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *CommonDao) GetTags(ids []int) ([]model.TinyTag, error) {
	ctxi := d.Context
	var tags []model.TinyTag
	err := d.db.Table(model.TableNameTag).Select("id,name").
		Where("id IN ?"+dbi.WithNotDeleted, ids).
		Find(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errcode.DBError, err, "GetTags")
	}
	return tags, nil
}

func (d *CommonDao) GetTagsByName(names []string) ([]model.TinyTag, error) {
	ctxi := d.Context
	var tags []model.TinyTag
	err := d.db.Table(model.TableNameTag).Select("id,name").
		Where("name IN ?"+dbi.WithNotDeleted, names).
		Find(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errcode.DBError, err, "GetTags")
	}
	return tags, nil
}
