package initialize

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

type RedisConfig struct {
	Addr        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

func (conf *RedisConfig) Generate() *redis.Pool {
	conf.IdleTimeout = conf.IdleTimeout * time.Second
	pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Addr)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	//closes = append(closes,pool.Close)
	return pool
}

func (init *Init) P2Redis() *redis.Pool {
	conf := &RedisConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}

	return conf.Generate()
}
