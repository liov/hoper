package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/scaffold/errcode"

	"github.com/liov/hoper/server/go/protobuf/content"
)

func (d *ContentDao) GetTopMoments(key string, pageNo int, PageSize int) ([]content.Moment, error) {
	ctxi := d.Context
	var moments []content.Moment
	exist, err := d.conn.Exists(ctxi.Base(), key).Result()
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.RedisErr, err, "GetTopMoments")
	}
	if exist == 0 {
		return nil, ctxi.RespErrorLog(errcode.DataLoss, err, "GetTopMoments")
	}
	d.conn.Pipelined(ctxi.Base(), func(pipe redis.Pipeliner) error {
		return nil
	})

	return moments, ctxi.RespErrorLog(errcode.RedisErr, err, "GetTopMoments")
}
