package redis

import (
	"github.com/redis/go-redis/v9"
)

type ContentDao struct {
	*redis.Client
}

func GetDao(client *redis.Client) *ContentDao {
	return &ContentDao{client}
}
