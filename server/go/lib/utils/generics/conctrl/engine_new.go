package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type TaskFunc[T any] func(context.Context) (T, error)

type TaskB[REQ, RES any] func(context.Context, REQ) (RES, error)

type TaskC[REQ, RES any] interface {
	TaskFunc[RES] | TaskB[REQ, RES]
}

type Task[T any] struct {
	Id        uint
	Kind      conctrl.Kind
	Do        TaskFunc[T]
	ErrHandle conctrl.ErrHandle
}

func NewTask[T any](taskFun TaskFunc[T]) *Task[T] {
	return &Task[T]{
		Do: taskFun,
	}
}

type Engine[T, V any] struct {
	ctrlEngine   *conctrl.Engine
	Properties   V
	resultChan   chan T
	resultHandle func(T)
	failChan     chan *Task[T]
	failHandle   func(*Task[T])
}

func NewEngine[T, V any](workerCount uint, props V, fun func(T), failFun func(*Task[T])) *Engine[T, V] {
	return &Engine[T, V]{
		ctrlEngine:   conctrl.NewEngine(workerCount),
		resultChan:   make(chan T),
		resultHandle: fun,
		failHandle:   failFun,
		failChan:     make(chan *Task[T]),
		Properties:   props,
	}
}

func (e *Engine[T, V]) Run(tasks ...*Task[T]) {
	go func() {
		for res := range e.resultChan {
			e.resultHandle(res)
		}
	}()
	go func() {
		for fail := range e.failChan {
			e.failHandle(fail)
		}
	}()
	ctrlTasks := make([]*conctrl.Task, 0, len(tasks))
	for _, task := range tasks {
		ctrlTasks = append(ctrlTasks, e.NewTask(task))
	}
	e.ctrlEngine.Run(ctrlTasks...)
}

func (e *Engine[T, V]) NewTask(task *Task[T]) *conctrl.Task {
	return &conctrl.Task{
		Do: func(ctx context.Context) error {
			res, err := task.Do(ctx)
			if err != nil {
				if task.ErrHandle != nil {
					task.ErrHandle(ctx, err)
				}
				e.failChan <- task
				return nil
			}
			e.resultChan <- res
			return nil
		},
	}
}
