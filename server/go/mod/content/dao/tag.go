package dao

import (
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"github.com/liov/hoper/server/go/mod/content/model"
	"github.com/liov/hoper/server/go/mod/protobuf/content"
)

const TagTableNameAlias = model.TagTableName + " a"

func (d *contentDao) GetContentTagDB(typ content.ContentType, refIds []uint64) ([]model.ContentTagRel, error) {
	ctxi := d.Ctx
	var tags []model.ContentTagRel
	err := d.Select("b.ref_id,a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.ContentTagTableName+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id IN (?) AND "+dbi.PostgreNotDeleted,
			typ, refIds).Find(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentTagDB")
	}
	return tags, nil
}

func (d *contentDao) GetTagsDB(names []string) ([]model.TinyTag, error) {
	ctxi := d.Ctx
	var tags []model.TinyTag
	err := d.Table(model.TagTableName).Select("id,name").
		Where("name IN (?) AND "+dbi.PostgreNotDeleted, names).
		Find(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetTagsDB")
	}
	return tags, nil
}

func (d *contentDao) GetTagsByRefIdDB(typ content.ContentType, refId uint64) ([]*content.TinyTag, error) {
	ctxi := d.Ctx
	var tags []*content.TinyTag
	err := d.Select("a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.ContentTagTableName+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id = ? AND "+dbi.PostgreNotDeleted,
			typ, refId).Scan(&tags).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetTagsByRefIdDB")
	}
	return tags, nil
}

func (d *contentDao) GetContentExtDB(typ content.ContentType, refIds []uint64) ([]*content.ContentExt, error) {
	ctxi := d.Ctx
	var exts []*content.ContentExt
	err := d.Table(model.ContentExtTableName).
		Where("type = ? AND ref_id IN (?)", typ, refIds).Find(&exts).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentTagDB")
	}
	return exts, nil
}
