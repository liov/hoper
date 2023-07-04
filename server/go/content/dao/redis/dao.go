package redis

import (
	"github.com/hopeio/zeta/context/http_context"
	"github.com/hopeio/zeta/utils/log"
)

type ContentRedisDao struct {
	*http_context.Context
	conn redis.Cmdable
}

func GetDao(ctx *http_context.Context, redis redis.Cmdable) *ContentRedisDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentRedisDao{ctx, redis}
}
