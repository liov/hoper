package conctrl

import (
	"context"
	"math/rand"
	"time"
)

type TimerTask struct {
	Times     uint
	FirstExec bool
	Do        BaseTaskFunc
}

func (task *TimerTask) Timer(ctx context.Context, interval time.Duration) {
	timer := time.NewTicker(interval)
	if task.FirstExec {
		task.Times = 1
		task.Do(ctx)
	}
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

func (task *TimerTask) RandTimer(ctx context.Context, start, end time.Duration) {
	range1 := end - start
	timer := time.NewTimer(time.Duration(rand.Intn(int(range1))) + start)
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
			timer.Reset(time.Duration(rand.Intn(int(range1))) + start)
		}
	}
}
