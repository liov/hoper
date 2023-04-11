package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/pandora/context/http_context"
	"github.com/liov/hoper/server/go/mod/content/dao/db"
	rdao "github.com/liov/hoper/server/go/mod/content/dao/redis"
	"gorm.io/gorm"
)

func GetDBDao(ctx *http_context.Context, d *gorm.DB) *db.ContentDBDao {
	return db.GetDao(ctx, d)
}

func GetRedisDao(ctx *http_context.Context, r redis.Cmdable) *rdao.ContentRedisDao {
	return rdao.GetDao(ctx, r)
}
