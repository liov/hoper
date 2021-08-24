package redisi

import (
	"context"
	"testing"

	goredis "github.com/go-redis/redis/v8"
)

const LoginUserKey = "LoginUserKey"

func BenchmarkGoRedis(b *testing.B) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     "192.168.1.204:6379",
		Password: "123456",
	})
	ctx := context.Background()
	conn := client.Conn(ctx)
	for i := 0; i < b.N; i++ {
		cmd := goredis.NewCmd(ctx, "HGETALL", LoginUserKey+"1")
		_ = conn.Process(ctx, cmd)
	}
}
