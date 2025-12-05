package data

import (
	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/initialize/dao/redis"
	"github.com/liov/hoper/server/go/content/data/db"
	rdao "github.com/liov/hoper/server/go/content/data/redis"
	"gorm.io/gorm"
)

func GetDBDao(ctx *httpctx.Context, d *gorm.DB) *db.ContentDao {
	return db.GetDao(ctx, d)
}

func GetRedisDao(ctx *httpctx.Context, r redis.Client) *rdao.ContentDao {
	return rdao.GetDao(ctx, r)
}
