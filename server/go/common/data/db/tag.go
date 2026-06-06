package db

import (
	"context"

	sqlx "github.com/hopeio/gox/database/sql"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/content/model"
	"go.uber.org/zap"
)

func (d *CommonDao) GetTags(ctx context.Context, ids []int) ([]model.TinyTag, error) {
	var tags []model.TinyTag
	err := d.Table(model.TableNameTag).Select("id,name").
		Where("id IN ?"+sqlx.WithNotDeleted, ids).
		Find(&tags).Error
	if err != nil {
		log.Errorw("GetTags faild", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	return tags, nil
}

func (d *CommonDao) GetTagsByName(ctx context.Context, names []string) ([]model.TinyTag, error) {
	var tags []model.TinyTag
	err := d.Table(model.TableNameTag).Select("id,name").
		Where("name IN ?"+sqlx.WithNotDeleted, names).
		Find(&tags).Error
	if err != nil {
		log.Errorw("GetTagsByName faild", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	return tags, nil
}
