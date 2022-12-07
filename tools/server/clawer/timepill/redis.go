package timepill

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type TimepillRedisDao struct {
	ctx   context.Context
	Redis *redis.Client
}
