package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/context/httpctx"
	redisi "github.com/hopeio/initialize/dao/redis"
	"github.com/hopeio/utils/log"
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
