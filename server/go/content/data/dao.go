package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/liov/hoper/server/go/content/data/db"
	rdao "github.com/liov/hoper/server/go/content/data/redis"
	"gorm.io/gorm"
)

func GetDBDao(ctx *httpctx.Context, d *gorm.DB) *db.ContentDao {
	return db.GetDao(ctx, d)
}

func GetRedisDao(ctx *httpctx.Context, r redis.Cmdable) *rdao.ContentDao {
	return rdao.GetDao(ctx, r)
}
