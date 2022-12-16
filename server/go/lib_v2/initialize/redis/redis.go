package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/server/go/lib_v2/initialize"
	"time"
)

type Config struct {
	Addr        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Index       int
}

func (conf *Config) Build() (*redis.Client, func() error) {
	conf.IdleTimeout = conf.IdleTimeout * time.Second
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		MinIdleConns: conf.MaxActive,
		IdleTimeout:  conf.IdleTimeout,
		DB:           conf.Index,
	})
	//closes = append(closes,pool.CloseDao)
	return client, func() error {
		return client.Close()
	}
}

type RedisDao = initialize.Dao[*Config, redis.Client]
