package redis

import (
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/hopeio/cherry/protobuf/errorcode"
	timei "github.com/hopeio/cherry/utils/time"
	"github.com/liov/hoper/server/go/content/confdao"
)

var limitErr = errorcode.TimeTooMuch.Message("您的操作过于频繁，请先休息一会儿。")

func (d *ContentRedisDao) Limit(l *confdao.Limit) error {
	ctxi := d
	ctx := ctxi.Context.Context()
	minuteKey := l.MinuteLimitKey + ctxi.AuthID
	dayKey := l.DayLimitKey + ctxi.AuthID

	var minuteIntCmd, dayIntCmd *redis.IntCmd
	_, err := d.conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteIntCmd = pipe.Incr(ctx, minuteKey)
		dayIntCmd = pipe.Incr(ctx, dayKey)
		return nil
	})
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "Incr")
	}

	if minuteIntCmd.Val() > l.MinuteLimitCount || dayIntCmd.Val() > l.DayLimitCount {
		return limitErr
	}
	var minuteDurationCmd, dayDurationCmd *redis.DurationCmd
	_, err = d.conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteDurationCmd = pipe.PTTL(ctx, minuteKey)
		dayDurationCmd = pipe.PTTL(ctx, dayKey)
		return nil
	})
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "PTTL")
	}

	_, err = d.conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if minuteDurationCmd.Val() < 0 {
			pipe.Expire(ctx, minuteKey, time.Minute)
		}
		if dayDurationCmd.Val() < 0 {
			pipe.Expire(ctx, dayKey, timei.TimeDay)
		}
		return nil
	})
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "Expire")
	}
	return nil
}
