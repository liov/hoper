package db

import (
	"github.com/actliboy/hoper/server/go/content/model"
	"github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/hopeio/lemon/protobuf/errorcode"
	dbi "github.com/hopeio/lemon/utils/dao/db/const"
)

const TagTableNameAlias = model.TagTableName + " a"

func (d *ContentDBDao) GetContentTag(typ content.ContentType, refIds []uint64) ([]model.ContentTagRel, error) {
	ctxi := d.Context
	var tags []model.ContentTagRel
	err := d.db.Select("b.ref_id,a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.ContentTagTableName+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id IN (?)"+dbi.WithNotDeleted,
			typ, refIds).Find(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentTag")
	}
	return tags, nil
}

func (d *ContentDBDao) GetTags(names []string) ([]model.TinyTag, error) {
	ctxi := d.Context
	var tags []model.TinyTag
	err := d.db.Table(model.TagTableName).Select("id,name").
		Where("name IN (?)"+dbi.WithNotDeleted, names).
		Find(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetTags")
	}
	return tags, nil
}

func (d *ContentDBDao) GetTagsByRefId(typ content.ContentType, refId uint64) ([]*content.TinyTag, error) {
	ctxi := d.Context
	var tags []*content.TinyTag
	err := d.db.Select("a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.ContentTagTableName+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id = ?"+dbi.WithNotDeleted,
			typ, refId).Scan(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetTagsByRefId")
	}
	return tags, nil
}

func (d *ContentDBDao) GetContentExt(typ content.ContentType, refIds []uint64) ([]*content.ContentExt, error) {
	ctxi := d.Context
	var exts []*content.ContentExt
	err := d.db.Table(model.ContentExtTableName).
		Where("type = ? AND ref_id IN (?)", typ, refIds).Find(&exts).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentTag")
	}
	return exts, nil
}
