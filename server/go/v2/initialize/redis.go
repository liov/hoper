package initialize

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/utils/reflect"
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

func (init *Init) P2Redis() *redis.Client {
	conf := &RedisConfig{}
	if exist := reflecti.GetFieldValue(init.conf, conf); !exist {
		return nil
	}

	return conf.Generate()
}
