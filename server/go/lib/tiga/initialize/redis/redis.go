package redis

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

func (conf *RedisConfig) generate() *redis.Client {
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

func (conf *RedisConfig) Generate() interface{} {
	return conf.generate()
}

type Redis struct {
	*redis.Client
	Conf RedisConfig
}

func (db *Redis) Config() interface{} {
	return &db.Conf
}

func (db *Redis) SetEntity(entity interface{}) {
	if client, ok := entity.(*redis.Client); ok {
		db.Client = client
	}
}

func (db *Redis) Close() error {
	return db.Client.Close()
}
