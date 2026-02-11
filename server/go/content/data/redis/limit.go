package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hopeio/scaffold/errcode"

	timex "github.com/hopeio/gox/time"
	"github.com/liov/hoper/server/go/global"
)

var limitErr = errcode.TimesTooMuch.Msg("您的操作过于频繁，请先休息一会儿。")

func (d *ContentDao) Limit(ctx context.Context, l *global.ContentLimit, userId uint64) error {

	userIdStr := strconv.FormatUint(userId, 10)
	minuteKey := l.MinuteLimitKey + userIdStr
	dayKey := l.DayLimitKey + userIdStr

	var minuteIntCmd, dayIntCmd *redis.IntCmd
	_, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteIntCmd = pipe.Incr(ctx, minuteKey)
		dayIntCmd = pipe.Incr(ctx, dayKey)
		return nil
	})
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}

	if minuteIntCmd.Val() > l.MinuteLimitCount || dayIntCmd.Val() > l.DayLimitCount {
		return limitErr
	}
	var minuteDurationCmd, dayDurationCmd *redis.DurationCmd
	_, err = d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		minuteDurationCmd = pipe.PTTL(ctx, minuteKey)
		dayDurationCmd = pipe.PTTL(ctx, dayKey)
		return nil
	})
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}

	_, err = d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if minuteDurationCmd.Val() < 0 {
			pipe.Expire(ctx, minuteKey, time.Minute)
		}
		if dayDurationCmd.Val() < 0 {
			pipe.Expire(ctx, dayKey, timex.Day)
		}
		return nil
	})
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}
