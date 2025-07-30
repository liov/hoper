package db

import (
	"github.com/hopeio/scaffold/errcode"
	dbi "github.com/hopeio/gox/datax/database/sql"
	commonmodel "github.com/liov/hoper/server/go/common/model"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/common"
	"github.com/liov/hoper/server/go/protobuf/content"
)

const TagTableNameAlias = commonmodel.TableNameTag + " a"

func (d *ContentDao) GetContentTag(typ content.ContentType, refIds []uint64) ([]model.ContentTagRel, error) {
	ctxi := d.Context
	var tags []model.ContentTagRel
	err := d.db.Select("b.ref_id,a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.TableNameContentTag+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id IN (?)"+dbi.WithNotDeleted,
			typ, refIds).Find(&tags).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "GetContentTag")
	}
	return tags, nil
}

func (d *ContentDao) GetTagsByRefId(typ content.ContentType, refId uint64) ([]*common.TinyTag, error) {
	ctxi := d.Context
	var tags []*common.TinyTag
	err := d.db.Select("a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.TableNameContentTag+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id = ?"+dbi.WithNotDeleted,
			typ, refId).Scan(&tags).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "GetTagsByRefId")
	}
	return tags, nil
}

func (d *ContentDao) GetStatistics(typ content.ContentType, refIds []uint64) ([]*content.Statistics, error) {
	ctxi := d.Context
	var exts []*content.Statistics
	err := d.db.Table(model.TableNameStatistics).
		Where("type = ? AND ref_id IN (?)", typ, refIds).Find(&exts).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "GetContentTag")
	}
	return exts, nil
}
