package conctrl

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/slices"
	"log"
	"sync"
	"time"
)

type Engine[KEY comparable, T, W any] struct {
	*BaseEngine[KEY, T, W]
	done        sync.Map
	TasksChan   chan []*Task[KEY, T]
	kindHandler []*KindHandler[KEY, T]
	errHandler  func(task *Task[KEY, T])
	errChan     chan *Task[KEY, T]
}

type KindHandler[KEY comparable, T any] struct {
	Skip   bool
	Ticker *time.Ticker
	// TODO 指定Kind的Handler
	HandleFun TaskFunc[KEY, T]
}

func NewEngine[KEY comparable, T, W any](workerCount uint) *Engine[KEY, T, W] {
	return &Engine[KEY, T, W]{
		BaseEngine: NewBaseEngine[KEY, T, W](workerCount),
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
	e.Release()
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
		if _, ok := e.done.Load(task.Key); ok {
			return nil
		}
	}
	return &BaseTask[KEY, T]{
		BaseTaskMeta: task.BaseTaskMeta,
		BaseTaskFunc: func(ctx context.Context) {
			if kindHandler != nil && kindHandler.Ticker != nil {
				<-kindHandler.Ticker.C
			}
			tasks, err := task.TaskFunc(ctx)
			if err != nil {
				task.ErrTimes++
				task.errs = append(task.errs, err)
				log.Println("执行失败", err)
				log.Println("重新执行,key :", task.Key)
				if task.ErrTimes < 5 {
					e.AsyncAddTask(task.Priority+1, task)
				}
				if task.ErrTimes == 5 {
					e.errChan <- task
				}
				return
			}
			if task.Key != zeroKey {
				e.done.Store(task.Key, struct{}{})
			}
			if len(tasks) > 0 {
				e.AsyncAddTask(task.Priority+1, tasks...)
			}
			return
		},
	}
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
	for _, kindHandler := range e.kindHandler {
		if kindHandler.Ticker != nil {
			kindHandler.Ticker.Stop()
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
