package conctrl

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/liov/hoper/server/go/lib/utils/slices"
	"golang.org/x/time/rate"
	"log"
	"time"
)

type Engine[KEY comparable, T, W any] struct {
	*BaseEngine[KEY, T, W]
	done *ristretto.Cache
	//TasksChan   chan []*Task[KEY, T]
	kindHandler []*KindHandler[KEY, T]
	errHandler  func(task *Task[KEY, T])
	errChan     chan *Task[KEY, T]
}

type KindHandler[KEY comparable, T any] struct {
	Skip    bool
	Ticker  *time.Ticker
	Limiter *rate.Limiter
	// TODO 指定Kind的Handler
	HandleFun TaskFunc[KEY, T]
}

func NewEngine[KEY comparable, T, W any](workerCount uint) *Engine[KEY, T, W] {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters:        1e4,   // number of keys to track frequency of (10M).
		MaxCost:            1e3,   // maximum cost of cache (MaxCost * 1MB).
		BufferItems:        64,    // number of keys per Get buffer.
		Metrics:            false, // number of keys per Get buffer.
		IgnoreInternalCost: true,
	})
	return &Engine[KEY, T, W]{
		BaseEngine: NewBaseEngine[KEY, T, W](workerCount),
		done:       cache,
		errHandler: func(task *Task[KEY, T]) {
			log.Println(task.errs)
		},
		errChan: make(chan *Task[KEY, T]),
	}
}

func (e *Engine[KEY, T, W]) SkipKind(kinds ...Kind) *Engine[KEY, T, W] {
	length := slices.Max(kinds) + 1
	if e.kindHandler == nil {
		e.kindHandler = make([]*KindHandler[KEY, T], length)
	}
	if int(length) > len(e.kindHandler) {
		e.kindHandler = append(e.kindHandler, make([]*KindHandler[KEY, T], int(length)-len(e.kindHandler))...)
	}
	for _, kind := range kinds {
		if e.kindHandler[kind] == nil {
			e.kindHandler[kind] = &KindHandler[KEY, T]{Skip: true}
		} else {
			e.kindHandler[kind].Skip = true
		}

	}
	return e
}
func (e *Engine[KEY, T, W]) StopAfter(interval time.Duration) *Engine[KEY, T, W] {
	time.AfterFunc(interval, e.Cancel)
	return e
}

func (e *Engine[KEY, T, W]) ErrHandler(errHandler func(task *Task[KEY, T])) *Engine[KEY, T, W] {
	e.errHandler = errHandler
	return e
}

func (e *Engine[KEY, T, W]) Timer(kind Kind, interval time.Duration) *Engine[KEY, T, W] {
	e.kindTimer(kind, time.NewTicker(interval))
	return e
}

func (e *Engine[KEY, T, W]) Limiter(kind Kind, r rate.Limit, b int) *Engine[KEY, T, W] {
	e.kindLimiter(kind, r, b)
	return e
}

// 多个kind共用一个timer
func (e *Engine[KEY, T, W]) KindGroupTimer(interval time.Duration, kinds ...Kind) *Engine[KEY, T, W] {
	ticker := time.NewTicker(interval)
	for _, kind := range kinds {
		e.kindTimer(kind, ticker)
	}
	return e
}

func (e *Engine[KEY, T, W]) kindTimer(kind Kind, ticker *time.Ticker) {
	if e.kindHandler == nil {
		e.kindHandler = make([]*KindHandler[KEY, T], int(kind)+1)
	}
	if int(kind)+1 > len(e.kindHandler) {
		e.kindHandler = append(e.kindHandler, make([]*KindHandler[KEY, T], int(kind)+1-len(e.kindHandler))...)
	}
	if e.kindHandler[kind] == nil {
		e.kindHandler[kind] = &KindHandler[KEY, T]{Ticker: ticker}
	} else {
		e.kindHandler[kind].Ticker = ticker
	}

}

func (e *Engine[KEY, T, W]) kindLimiter(kind Kind, r rate.Limit, b int) {
	if e.kindHandler == nil {
		e.kindHandler = make([]*KindHandler[KEY, T], int(kind)+1)
	}
	if int(kind)+1 > len(e.kindHandler) {
		e.kindHandler = append(e.kindHandler, make([]*KindHandler[KEY, T], int(kind)+1-len(e.kindHandler))...)
	}
	if e.kindHandler[kind] == nil {
		e.kindHandler[kind] = &KindHandler[KEY, T]{Limiter: rate.NewLimiter(r, b)}
	} else {
		e.kindHandler[kind].Limiter = rate.NewLimiter(r, b)
	}

}

func (e *Engine[KEY, T, W]) Run(tasks ...*Task[KEY, T]) {
	baseTasks := make([]*BaseTask[KEY, T], 0, len(tasks))
	for _, task := range tasks {
		baseTasks = append(baseTasks, e.BaseTask(task))
	}
	go func() {
		for group := range e.errChan {
			e.errHandler(group)
		}
	}()
	e.BaseEngine.Run(baseTasks...)
}

func (e *Engine[KEY, T, W]) ReRun(tasks ...*Task[KEY, T]) {
	if e.isRunning {
		e.AddTasks(0, tasks...)
		return
	}

	baseTasks := make([]*BaseTask[KEY, T], 0, len(tasks))
	for _, task := range tasks {
		baseTasks = append(baseTasks, e.BaseTask(task))
	}
	e.BaseEngine.Run(baseTasks...)
}

func (e *Engine[KEY, T, W]) BaseTask(task *Task[KEY, T]) *BaseTask[KEY, T] {

	if task == nil {
		return nil
	}

	var kindHandler *KindHandler[KEY, T]
	if e.kindHandler != nil && int(task.Kind) < len(e.kindHandler) {
		kindHandler = e.kindHandler[task.Kind]
	}

	if kindHandler != nil && kindHandler.Skip {
		return nil
	}

	zeroKey := *new(KEY)

	if task.Key != zeroKey {
		if _, ok := e.done.Get(task.Key); ok {
			return nil
		}
	}
	return &BaseTask[KEY, T]{
		BaseTaskMeta: task.BaseTaskMeta,
		BaseTaskFunc: func(ctx context.Context) {
			if kindHandler != nil {
				if kindHandler.Ticker != nil {
					<-kindHandler.Ticker.C
				}
				if kindHandler.Limiter != nil {
					kindHandler.Limiter.Wait(ctx)
				}
			}
			tasks, err := task.TaskFunc(ctx)
			if err != nil {
				task.ErrTimes++
				task.errs = append(task.errs, err)
				log.Println(task.Key, "执行失败", err, "重新执行")
				if task.ErrTimes < 5 {
					e.AsyncAddTask(task.Priority+1, task)
				}
				if task.ErrTimes == 5 {
					e.errChan <- task
				}
				return
			}
			if task.Key != zeroKey {
				e.done.SetWithTTL(task.Key, struct{}{}, 1, time.Hour)
			}
			if len(tasks) > 0 {
				e.AsyncAddTask(task.Priority+1, tasks...)
			}
			return
		},
	}
}

func (e *Engine[KEY, T, W]) AddTask(task *Task[KEY, T]) {
	e.BaseEngine.AddTask(e.BaseTask(task))
}

// IDEA的bug,可以编译过去，但是报错，提示方法是嵌套结构体的方法,为了跳转写的包装方法
func (e *Engine[KEY, T, W]) AddTasksWrapper(generation int, tasks ...*Task[KEY, T]) {
	e.AddTasks(generation, tasks...)
}

func (e *Engine[KEY, T, W]) AddTasks(generation int, tasks ...*Task[KEY, T]) {
	for _, req := range tasks {
		if req != nil {
			req.Priority += generation
			e.BaseEngine.AddTask(e.BaseTask(req))
		}
	}
}

func (e *Engine[KEY, T, W]) AsyncAddTask(generation int, tasks ...*Task[KEY, T]) {
	go func() {
		for _, task := range tasks {
			if task != nil {
				task.Priority += generation
				e.BaseEngine.AddTask(e.BaseTask(task))
			}
		}
	}()
}

func (e *Engine[KEY, T, W]) AddFixedTask(workerId int, task *Task[KEY, T]) {
	e.BaseEngine.AddFixedTask(workerId, task.BaseTask(func(tasks []*Task[KEY, T], err error) {
		if err != nil {
			e.AddFixedTask(workerId, task)
		} else {
			for _, task := range tasks {
				e.AddFixedTask(workerId, task)
			}
		}
	}))
}

func (e *Engine[KEY, T, W]) RunSingleWorker(tasks ...*Task[KEY, T]) {
	e.limitWorkerCount = 1
	e.Run(tasks...)
}

func (e *Engine[KEY, T, W]) Release() {
	e.BaseEngine.Release()
	e.done.Close()
	for _, kindHandler := range e.kindHandler {
		if kindHandler != nil {
			if kindHandler.Ticker != nil {
				kindHandler.Ticker.Stop()
			}
			if kindHandler.Limiter != nil {
				kindHandler.Limiter = nil
			}
		}
	}
}

func NewTask[KEY comparable, T any](baseTask *BaseTask[KEY, T]) *Task[KEY, T] {
	return &Task[KEY, T]{
		TaskMeta: TaskMeta[KEY]{BaseTaskMeta: baseTask.BaseTaskMeta},
		TaskFunc: func(ctx context.Context) ([]*Task[KEY, T], error) {
			baseTask.BaseTaskFunc(ctx)
			return nil, nil
		},
	}
}

func AnonymousTask[KEY comparable, T any](fun BaseTaskFunc) *Task[KEY, T] {
	return &Task[KEY, T]{
		TaskFunc: func(ctx context.Context) ([]*Task[KEY, T], error) {
			fun(ctx)
			return nil, nil
		},
	}
}
