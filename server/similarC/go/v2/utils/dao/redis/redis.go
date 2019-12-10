package redis

import "github.com/gomodule/redigo/redis"

func Do(pool *redis.Pool,commandName string, args ...interface{} ) (interface{},error) {
	conn := pool.Get()
	defer conn.Close()
	return conn.Do(commandName,args...)
}
