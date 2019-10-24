package initialize

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P2Redis() *redis.Pool {
	conf :=RedisConfig{}
	if exist := reflect3.GetFieldValue(i.conf,&conf);!exist{
		return nil
	}
	return &redis.Pool{
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
}
