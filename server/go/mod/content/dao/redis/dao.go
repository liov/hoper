package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/server/go/lib/context/http_context"
	"github.com/liov/hoper/server/go/lib/utils/log"
)

type ContentRedisDao struct {
	*http_context.Ctx
	conn redis.Cmdable
}

func GetDao(ctx *http_context.Ctx, redis redis.Cmdable) *ContentRedisDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentRedisDao{ctx, redis}
}
