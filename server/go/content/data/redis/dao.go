package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/utils/log"
)

type ContentRedisDao struct {
	*httpctx.Context
	conn redis.Cmdable
}

func GetDao(ctx *httpctx.Context, redis redis.Cmdable) *ContentRedisDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentRedisDao{ctx, redis}
}
