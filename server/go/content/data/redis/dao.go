package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hopeio/gox/log"
	redisi "github.com/hopeio/initialize/dao/redis"
)

type ContentDao struct {
	*redis.Client
}

func GetDao(ctx context.Context, redis redisi.Client) *ContentDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDao{redis.WithContext(ctx)}
}
