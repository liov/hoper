package redis

import (
	"github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/pandora/protobuf/errorcode"
)

func (d *ContentRedisDao) GetTopMoments(key string, pageNo int, PageSize int) ([]content.Moment, error) {
	ctxi := d.Context
	var moments []content.Moment
	exist, err := d.conn.Exists(ctxi.Context, key).Result()
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "GetTopMoments")
	}
	if exist == 0 {
		return nil, ctxi.ErrorLog(errorcode.DataLoss, err, "GetTopMoments")
	}
	d.conn.Pipelined(ctxi.Context, func(pipe redis.Pipeliner) error {
		return nil
	})

	return moments, ctxi.ErrorLog(errorcode.RedisErr, err, "GetTopMoments")
}
