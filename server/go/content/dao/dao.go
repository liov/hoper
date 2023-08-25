package dao

import (
	"github.com/actliboy/hoper/server/go/content/dao/db"
	rdao "github.com/actliboy/hoper/server/go/content/dao/redis"
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/lemon/context/http_context"
	"gorm.io/gorm"
)

func GetDBDao(ctx *http_context.Context, d *gorm.DB) *db.ContentDBDao {
	return db.GetDao(ctx, d)
}

func GetRedisDao(ctx *http_context.Context, r redis.Cmdable) *rdao.ContentRedisDao {
	return rdao.GetDao(ctx, r)
}
