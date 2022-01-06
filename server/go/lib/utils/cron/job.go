package cron

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type RedisTo struct {
}

func (RedisTo) Run() {
	log.Info("定时任务执行")
}
