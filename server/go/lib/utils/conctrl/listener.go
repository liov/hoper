package conctrl

import (
	"context"
	"time"
)

func Timer(ctx context.Context, task Task) {
	timer := time.NewTicker(time.Second)
	task.Do(ctx)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			task.Do(ctx)
		}
	}
}
