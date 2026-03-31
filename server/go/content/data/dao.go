package data

import (
	"github.com/liov/hoper/server/go/content/data/db"
	rdao "github.com/liov/hoper/server/go/content/data/redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetDBDao(d *gorm.DB) *db.ContentDao {
	return db.GetDao(d)
}

func GetRedisDao(client *redis.Client) *rdao.ContentDao {
	return rdao.GetDao(client)
}
