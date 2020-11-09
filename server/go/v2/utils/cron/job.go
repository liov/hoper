package cron

import (
	"github.com/liov/hoper/go/v2/utils/log"
)

type RedisTo struct {
}

func (RedisTo) Run() {
	log.Info("定时任务执行")
}
