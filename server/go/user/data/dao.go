package data

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/server/go/user/data/db"
	rdao "github.com/liov/hoper/server/go/user/data/redis"
	"gorm.io/gorm"
)

func GetDBDao(ctx context.Context, d *gorm.DB) *db.UserDao {
	return db.GetUserDao(ctx, d)
}

func GetRedisDao(ctx context.Context, c *redis.Client) *rdao.UserDao {
	return rdao.GetUserDao(ctx, c)
}
