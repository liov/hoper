package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/utils/log"
)

type ContentDao struct {
	*httpctx.Context
	conn redis.Cmdable
}

func GetDao(ctx *httpctx.Context, redis redis.Cmdable) *ContentDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDao{ctx, redis}
}
