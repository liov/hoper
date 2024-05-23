package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/liov/hoper/server/go/user/data/db"
	rdao "github.com/liov/hoper/server/go/user/data/redis"
	"gorm.io/gorm"
)

func GetDBDao(ctx *httpctx.Context, d *gorm.DB) *db.UserDao {
	return db.GetUserDao(ctx, d)
}

func GetRedisDao(ctx *httpctx.Context, c *redis.Client) *rdao.UserDao {
	return rdao.GetUserDao(ctx, c)
}
