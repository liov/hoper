package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/server/go/user/data/db"
	rdao "github.com/liov/hoper/server/go/user/data/redis"
	"gorm.io/gorm"
)

func GetDBDao(d *gorm.DB) *db.UserDao {
	return db.GetUserDao(d)
}

func GetRedisDao(c *redis.Client) *rdao.UserDao {
	return rdao.GetUserDao(c)
}
