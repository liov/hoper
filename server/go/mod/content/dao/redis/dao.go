package redis

import (
	contexti "github.com/actliboy/hoper/server/go/lib/tiga/context"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/go-redis/redis/v8"
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
