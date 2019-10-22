package initialize

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P2Redis() {
	redisConf:=RedisConfig{}
	if exist := reflect3.GetFieldValue(i.conf,&redisConf);!exist{
		return
	}
	url := fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
	reflect3.SetFieldValue(i.dao,&redis.Pool{
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
	})
}
