package service

import (
	"fmt"
	"time"

	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"github.com/liov/hoper/go/v2/utils/log"
)

func Limit(ctxi *user.Ctx, minuteLimit string, minuteLimitCount int64, dayLimit string, dayLimitCount int64, userID uint64) error {
	ctx := ctxi.Context
	conn := dao.Dao.Redis
	minuteKey := fmt.Sprintf("%s%d", minuteLimit, userID)
	minuteCount, minuteErr := redisi.Int64(conn.Get(ctx, minuteKey).Result())

	if minuteErr == nil && minuteCount >= minuteLimitCount {
		return errorcode.TimeTooMuch.Message("您的操作过于频繁，请先休息一会儿。")
	}

	minuteRemainingTime, _ := conn.TTL(ctx, minuteKey).Result()
	if minuteRemainingTime < 0 || minuteRemainingTime > 60*time.Second {
		minuteRemainingTime = 60 * time.Second
	}

	if err := conn.SetEX(ctx, minuteKey, minuteCount+1, minuteRemainingTime).Err(); err != nil {
		return errorcode.SysError
	}

	dayKey := fmt.Sprintf("%s%d", dayLimit, userID)
	dayCount, dayErr := redisi.Int64(conn.Get(ctx, dayKey).Result())
	if dayErr == nil && dayCount >= dayLimitCount {
		return errorcode.TimeTooMuch.Message("您今天的操作过于频繁，请先休息一会儿。")
	}

	dayRemainingTime, _ := conn.TTL(ctx, dayKey).Result()
	secondsOfDay := 24 * 60 * 60 * time.Second
	if dayRemainingTime < 0 || dayRemainingTime > secondsOfDay {
		dayRemainingTime = secondsOfDay
	}

	if err := conn.SetEX(ctx, dayKey, dayCount+1, dayRemainingTime).Err(); err != nil {
		log.Error("redis set failed:", err)
		return errorcode.SysError
	}
	return nil
}
