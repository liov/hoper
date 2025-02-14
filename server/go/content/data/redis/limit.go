package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/scaffold/errcode"
	"time"

	timei "github.com/hopeio/utils/time"
	"github.com/liov/hoper/server/go/content/global"
)

var limitErr = errcode.TimesTooMuch.Msg("您的操作过于频繁，请先休息一会儿。")

func (d *ContentDao) Limit(l *global.Limit) error {
	ctxi := d
	ctx := ctxi.Base()
	minuteKey := l.MinuteLimitKey + ctxi.AuthID
	dayKey := l.DayLimitKey + ctxi.AuthID

	var minuteIntCmd, dayIntCmd *redis.IntCmd
	_, err := d.conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteIntCmd = pipe.Incr(ctx, minuteKey)
		dayIntCmd = pipe.Incr(ctx, dayKey)
		return nil
	})
	if err != nil {
		return ctxi.RespErrorLog(errcode.RedisErr, err, "Incr")
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
		return ctxi.RespErrorLog(errcode.RedisErr, err, "PTTL")
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
		return ctxi.RespErrorLog(errcode.RedisErr, err, "Expire")
	}
	return nil
}
