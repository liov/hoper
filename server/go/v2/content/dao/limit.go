package dao

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
)


var limitErr = errorcode.TimeTooMuch.Message("您的操作过于频繁，请先休息一会儿。")

func (d *dao) Limit(ctxi *user.Ctx,l *conf.Limit) error {
	ctx := ctxi.Context
	conn := d.Redis
	minuteKey := l.MinuteLimit + ctxi.IdStr
	dayKey := l.DayLimit + ctxi.IdStr

	var minuteIntCmd, dayIntCmd *redis.IntCmd
	_, err := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteIntCmd = pipe.Incr(ctx, minuteKey)
		dayIntCmd = pipe.Incr(ctx, dayKey)
		return nil
	})
	if err != nil {
		ctxi.Error(err.Error())
		return errorcode.SysError
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
		ctxi.Error(err.Error())
		return errorcode.SysError
	}

	_, err = conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if minuteDurationCmd.Val() < 0 {
			pipe.Expire(ctx, minuteKey, time.Minute)
		}
		if dayDurationCmd.Val() < 0 {
			pipe.Expire(ctx, dayKey, redisi.TimeDay)
		}
		return nil
	})
	if err != nil {
		ctxi.Error(err.Error())
		return errorcode.SysError
	}
	return nil
}




