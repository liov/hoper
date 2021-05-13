package cron

import (
	"github.com/liov/hoper/v2/utils/log"
)

type RedisTo struct {
}

func (RedisTo) Run() {
	log.Info("定时任务执行")
}
