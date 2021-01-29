package redisi

import (
	"github.com/go-redis/redis/v8"
)

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