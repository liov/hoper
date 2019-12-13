package runtime2

import (
	"context"

	"go.uber.org/atomic"
)

//监测子协程是否跑完，需要代码层面的配合
//一点都不优雅
type NumGoroutine struct {
	context.Context
	context.CancelFunc
	start, end *atomic.Int32
	finished   bool
}

func New() *NumGoroutine {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &NumGoroutine{
		Context:    ctx,
		CancelFunc: cancel,
		start:      atomic.NewInt32(0),
		end:        atomic.NewInt32(0),
	}
}

func (ng *NumGoroutine) Start() {
	ng.start.Add(1)
}

func (ng *NumGoroutine) End(fn func()) {
	ng.end.Add(1)
	if ng.start.Load() == ng.end.Load() {
		fn()
		ng.finished = true
	}
}

func (ng *NumGoroutine) IsFinished() bool {
	return ng.finished
}

//有没有可能父协程return后，子协程还没开始执行
func (ng *NumGoroutine) Monitor() func(fn func()) {
	ng.Start()
	return ng.End
}

func (ng *NumGoroutine) Cancel() {
	ng.CancelFunc()
}
