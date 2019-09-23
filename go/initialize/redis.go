package initialize

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/liov/hoper/go/user/internal/config"
)

func (i *Init) P2Redis() *redis.Pool {
	var redisConf = &config.Config.Redis
	url := fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
	return &redis.Pool{
		MaxIdle:     redisConf.MaxIdle,
		MaxActive:   redisConf.MaxActive,
		IdleTimeout: redisConf.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}
			if redisConf.Password != "" {
				if _, err := c.Do("AUTH", redisConf.Password); err != nil {
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
