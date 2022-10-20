package conctrl

import (
	"context"
	"time"
)

type TimerTask struct {
	Times uint
	Do    TaskFunc
}

func Timer(ctx context.Context, task *TimerTask, interval time.Duration) {
	timer := time.NewTicker(interval)
	task.Times = 1
	task.Do(ctx)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			task.Times++
			task.Do(ctx)
		}
	}
}
