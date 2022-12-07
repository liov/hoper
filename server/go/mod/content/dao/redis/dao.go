package redis

import (
	"github.com/go-redis/redis/v8"
	contexti "github.com/liov/hoper/server/go/lib/context"
	"github.com/liov/hoper/server/go/lib/utils/log"
)

type ContentRedisDao struct {
	*contexti.Ctx
	conn redis.Cmdable
}

func GetDao(ctx *contexti.Ctx, redis redis.Cmdable) *ContentRedisDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentRedisDao{ctx, redis}
}
