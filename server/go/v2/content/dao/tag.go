package dao

import (
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	"gorm.io/gorm"
)

func (d *contentDao) GetTagContentDB(db *gorm.DB,typ content.ContentType,refIds []uint64) ([]model.TagContent,error) {
	var tags []model.TagContent
	err := db.Select("b.ref_id,a.id,a.name").Table("tag a").
		Joins(`LEFT JOIN content_tag b ON a.Id = b.tag_id`).
		Where("b.type = ? AND b.ref_id IN (?) AND deleted_at = ?",
		typ,refIds, dbi.PostgreZeroTime).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags,nil
}

func (d *contentDao) GetTagsDB(db *gorm.DB,names []string) ([]model.TinyTag,error) {
	var tags []model.TinyTag
	err:= db.Table("tag").Select("id,name").
		Where("name IN (?) AND deleted_at = ?", names,dbi.PostgreZeroTime).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags,nil
}

func (d *contentDao) GetTagsByRefIdDB(db *gorm.DB,typ content.ContentType,refId uint64) ([]*content.TinyTag,error) {
	var tags []*content.TinyTag
	err := db.Select("a.id,a.name").Table("tag a").
		Joins(`LEFT JOIN content_tag b ON a.Id = b.tag_id`).
		Where("b.type = ? AND b.ref_id = ? AND deleted_at = ?",
			typ,refId, dbi.PostgreZeroTime).Scan(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags,nil
}
