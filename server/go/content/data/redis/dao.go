package redis

import (
	"github.com/go-redis/redis/v8"
)

type ContentDao struct {
	*redis.Client
}

func GetDao(client *redis.Client) *ContentDao {
	return &ContentDao{client}
}
