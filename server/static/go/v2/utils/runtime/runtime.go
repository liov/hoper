package runtime2

import (
	"context"

	"go.uber.org/atomic"
)

//监测子协程是否跑完，需要代码层面的配合
//一点都不优雅
//没什么优不优雅，这就像wg.Add(1)
type NumGoroutine struct {
	context.Context
	context.CancelFunc
	start, end *atomic.Int32
	callback   func()
}

func New(callback func()) *NumGoroutine {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &NumGoroutine{
		Context:    ctx,
		CancelFunc: cancel,
		start:      atomic.NewInt32(0),
		end:        atomic.NewInt32(0),
		callback:   callback,
	}
}

func (ng *NumGoroutine) Start() {
	ng.start.Add(1)
}

func (ng *NumGoroutine) End() {
	ng.end.Add(1)
	if ng.start.Load() == ng.end.Load() {
		ng.callback()
	}
}

//有没有可能父协程return后，子协程还没开始执行
func (ng *NumGoroutine) Monitor(fn func()) {
	defer ng.End()
	ng.Start()
	fn()
}

func (ng *NumGoroutine) Cancel() {
	ng.CancelFunc()
}
