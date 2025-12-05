package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/log"
	redisi "github.com/hopeio/initialize/dao/redis"
)

type ContentDao struct {
	*httpctx.Context
	conn redis.Cmdable
}

func GetDao(ctx *httpctx.Context, redis redisi.Client) *ContentDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDao{ctx, redis.Client}
}
