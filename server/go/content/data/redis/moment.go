package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/protobuf/errcode"
	"github.com/liov/hoper/server/go/protobuf/content"
)

func (d *ContentDao) GetTopMoments(key string, pageNo int, PageSize int) ([]content.Moment, error) {
	ctxi := d.Context
	var moments []content.Moment
	exist, err := d.conn.Exists(ctxi.BaseContext(), key).Result()
	if err != nil {
		return nil, ctxi.ErrorLog(errcode.RedisErr, err, "GetTopMoments")
	}
	if exist == 0 {
		return nil, ctxi.ErrorLog(errcode.DataLoss, err, "GetTopMoments")
	}
	d.conn.Pipelined(ctxi.BaseContext(), func(pipe redis.Pipeliner) error {
		return nil
	})

	return moments, ctxi.ErrorLog(errcode.RedisErr, err, "GetTopMoments")
}
