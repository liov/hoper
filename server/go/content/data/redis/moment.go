package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hopeio/scaffold/errcode"

	"github.com/liov/hoper/server/go/protobuf/content"
)

func (d *ContentDao) GetTopMoments(ctx context.Context, key string, pageNo int, PageSize int) ([]content.Moment, error) {
	var moments []content.Moment
	exist, err := d.Exists(ctx, key).Result()
	if err != nil {
		return nil, errcode.RedisErr.Wrap(err)
	}
	if exist == 0 {
		return nil, errcode.DataLoss
	}
	d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		return nil
	})

	return moments, nil
}
