package db

import (
	"context"

	"github.com/hopeio/scaffold/errcode"
	commonmodel "github.com/liov/hoper/server/go/common/model"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/common"
	"github.com/liov/hoper/server/go/protobuf/content"
)

const TagTableNameAlias = commonmodel.TableNameTag + " a"

func (d *ContentDao) GetContentTag(ctx context.Context, typ content.ContentType, refIds []uint64) ([]model.ContentTagRel, error) {
	var tags []model.ContentTagRel
	err := d.Select("b.ref_id,a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.TableNameContentTag+` b ON a.id = b.tag_id`).
		Where("b.type = ? AND b.ref_id IN (?) AND b.deleted_at IS NULL",
			typ, refIds).Find(&tags).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return tags, nil
}

func (d *ContentDao) GetTagsByRefId(ctx context.Context, refId uint64) ([]*common.TinyTag, error) {
	var tags []*common.TinyTag
	err := d.Select("a.id,a.name").Table(TagTableNameAlias).
		Joins(`LEFT JOIN `+model.TableNameContentTag+` b ON a.id = b.tag_id`).
		Where("b.ref_id = ? AND b.deleted_at IS NULL",
			refId).Scan(&tags).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return tags, nil
}

func (d *ContentDao) GetStatistics(ctx context.Context, typ content.ContentType, refIds []uint64) ([]*content.Statistics, error) {
	var exts []*content.Statistics
	err := d.Table(model.TableNameStatistics).
		Where("type = ? AND ref_id IN (?)", typ, refIds).Find(&exts).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return exts, nil
}
