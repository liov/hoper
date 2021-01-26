package redisi

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func Do(ctx context.Context, conn *redis.Conn, args ...interface{}) (interface{}, error) {
	cmd := redis.NewCmd(ctx, args...)
	_ = conn.Process(ctx, cmd)
	return cmd.Result()
}

type Cmds []redis.Cmder

func (l Cmds) Last() redis.Cmder {
	return l[len(l)-1]
}

const(
	DEL = "DEL"
	HGETALL = "HGETALL"
	HMSET = "HMSET"
	SET = "SET"
)