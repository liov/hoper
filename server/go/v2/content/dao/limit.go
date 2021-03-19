package dao

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	timei "github.com/liov/hoper/go/v2/utils/time"
)


var limitErr = errorcode.TimeTooMuch.Message("您的操作过于频繁，请先休息一会儿。")

func (d *contentDao) LimitRedis(conn redis.Cmdable,l *conf.Limit) error {
	ctxi:=d.Ctx
	ctx := ctxi.Context
	minuteKey := l.MinuteLimit + ctxi.IdStr
	dayKey := l.DayLimit + ctxi.IdStr

	var minuteIntCmd, dayIntCmd *redis.IntCmd
	_, err := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteIntCmd = pipe.Incr(ctx, minuteKey)
		dayIntCmd = pipe.Incr(ctx, dayKey)
		return nil
	})
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err,"Incr")
	}

	if minuteIntCmd.Val() > l.MinuteLimitCount || dayIntCmd.Val() > l.DayLimitCount {
		return limitErr
	}
	var minuteDurationCmd, dayDurationCmd *redis.DurationCmd
	_, err = conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteDurationCmd = pipe.PTTL(ctx, minuteKey)
		dayDurationCmd = pipe.PTTL(ctx, dayKey)
		return nil
	})
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err,"PTTL")
	}

	_, err = conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if minuteDurationCmd.Val() < 0 {
			pipe.Expire(ctx, minuteKey, time.Minute)
		}
		if dayDurationCmd.Val() < 0 {
			pipe.Expire(ctx, dayKey, timei.TimeDay)
		}
		return nil
	})
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err,"Expire")
	}
	return nil
}




