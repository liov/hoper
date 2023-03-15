package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Addr        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Index       int
}

func (conf *Config) Build() *redis.Client {
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

type Redis struct {
	*redis.Client
	Conf Config
}

func (db *Redis) Config() any {
	return &db.Conf
}

func (db *Redis) SetEntity() {
	db.Client = db.Conf.Build()
}

func (db *Redis) Close() error {
	return db.Client.Close()
}
