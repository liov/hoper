package redisi

import (
	"context"
	"testing"

	goredis "github.com/go-redis/redis/v8"
	modelconst "github.com/liov/hoper/v2/user/model"
)

/*func BenchmarkRedigo(b *testing.B) {
	pool := &redis.Pool{
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: 200,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "192.168.1.204:6379")
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("AUTH", "123456"); err != nil {
				c.Close()
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	conn:=pool.Get()
	for i:=0;i<b.N;i++{
		conn.Do("HGETALL", modelconst.LoginUserKey + "1")
	}
}
*/
func BenchmarkGoRedis(b *testing.B) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     "192.168.1.204:6379",
		Password: "123456",
	})
	ctx := context.Background()
	conn := client.Conn(ctx)
	for i := 0; i < b.N; i++ {
		cmd := goredis.NewCmd(ctx, "HGETALL", modelconst.LoginUserKey+"1")
		_ = conn.Process(ctx, cmd)
	}
}
