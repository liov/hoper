package initialize

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Index       int
}

func (conf *RedisConfig) Generate() *redis.Client {
	conf.IdleTimeout = conf.IdleTimeout * time.Second
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		MinIdleConns: conf.MaxActive,
		IdleTimeout:  conf.IdleTimeout,
		DB:           conf.Index,
	})
	//closes = append(closes,pool.CloseDao)
	return client
}

func (init *Inject) P2Redis() *redis.Client {
	conf := &RedisConfig{}
	if exist := init.SetConf(conf); !exist {
		return nil
	}

	return conf.Generate()
}
