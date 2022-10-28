package v2

import (
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"sync"
	"time"
)

type Engine[T TaskProps] struct {
	ctrlEngine   *conctrl.BaseEngine
	visited      sync.Map
	reqsChan     chan []*Task[T]
	excludeKinds []bool
	timer        []*time.Ticker
}

func New[A any, T TaskProps](workerCount uint) *Engine[T] {
	return &Engine[T]{
		reqsChan:   make(chan []*Task[T]),
		ctrlEngine: conctrl.NewEngine(workerCount),
	}
}

func (e *Engine[T]) SkipKind(kinds ...conctrl.Kind) *Engine[T] {
	length := slices.Max(kinds) + 1
	if e.excludeKinds == nil {
		e.excludeKinds = make([]bool, length)
	}
	if int(length) > len(e.excludeKinds) {
		e.excludeKinds = append(e.excludeKinds, make([]bool, int(length)-len(e.excludeKinds))...)
	}
	for _, kind := range kinds {
		e.excludeKinds[kind] = true
	}
	return e
}

func (e *Engine[T]) Timer(kind conctrl.Kind, interval time.Duration) *Engine[T] {
	if e.timer == nil {
		e.timer = make([]*time.Ticker, int(kind)+1)
	}
	if int(kind)+1 > len(e.timer) {
		e.timer = append(e.timer, make([]*time.Ticker, int(kind)+1-len(e.timer))...)
	}
	e.timer[kind] = time.NewTicker(interval)
	return e
}

func (e *Engine[T]) Run(reqs ...*Task[T]) {

	go func() {
		for reqs := range e.reqsChan {
			for _, req := range reqs {
				if e.excludeKinds != nil && int(req.Kind) < len(e.excludeKinds) && e.excludeKinds[req.Kind] {
					continue
				}
				e.ctrlEngine.AddTask(e.NewTask(req))
			}
		}
	}()
	tasks := make([]*conctrl.Task, 0, len(reqs))
	for _, req := range reqs {
		tasks = append(tasks, e.NewTask(req))
	}
	e.ctrlEngine.Run(tasks...)
}

func (e *Engine[T]) NewTask(req *Task[T]) *conctrl.Task {
	fun := req.Props.NewTaskFun(req.Id, req.Kind)
	if fun == nil {
		return nil
	}
	return &conctrl.Task{
		Id:   req.Id,
		Kind: req.Kind,
		Do:   fun,
	}
}
