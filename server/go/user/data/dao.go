package data

import (
	"github.com/liov/hoper/server/go/user/data/db"
	rdao "github.com/liov/hoper/server/go/user/data/redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetDBDao(d *gorm.DB) *db.UserDao {
	return db.GetUserDao(d)
}

func GetRedisDao(c *redis.Client) *rdao.UserDao {
	return rdao.GetUserDao(c)
}
