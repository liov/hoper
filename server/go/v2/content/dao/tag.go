package dao

import (
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	"gorm.io/gorm"
)

const TagTableNameAlias = model.TagTableName + " a"

func (d *contentDao) GetContentTagDB(db *gorm.DB, typ content.ContentType, refIds []uint64) ([]model.ContentTagRel, error) {
	var tags []model.ContentTagRel
	err := db.Select("b.ref_id,a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.ContentTagTableName+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id IN (?) AND "+dbi.PostgreNotDeleted,
			typ, refIds).Find(&tags).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError, err, "GetContentTagDB")
	}
	return tags, nil
}

func (d *contentDao) GetTagsDB(db *gorm.DB, names []string) ([]model.TinyTag, error) {
	var tags []model.TinyTag
	err := db.Table(model.TagTableName).Select("id,name").
		Where("name IN (?) AND "+dbi.PostgreNotDeleted, names).
		Find(&tags).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError, err, "GetTagsDB")
	}
	return tags, nil
}

func (d *contentDao) GetTagsByRefIdDB(db *gorm.DB, typ content.ContentType, refId uint64) ([]*content.TinyTag, error) {
	var tags []*content.TinyTag
	err := db.Select("a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.ContentTagTableName+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id = ? AND "+dbi.PostgreNotDeleted,
			typ, refId).Scan(&tags).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError, err, "GetTagsByRefIdDB")
	}
	return tags, nil
}
