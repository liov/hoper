package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type TaskFuncA func(context.Context)

type TaskFuncB func(context.Context) error

type TaskFuncC[T any] func(context.Context) (T, error)

type TaskFuncD[REQ, RES any] func(context.Context, REQ) (RES, error)

type TaskFunc[REQ, RES any] interface {
	TaskFuncA | TaskFuncB | TaskFuncC[RES] | TaskFuncD[REQ, RES]
}

type Task[REQ, RES any, FUN TaskFunc[REQ, RES]] struct {
	Id   uint
	Kind conctrl.Kind
	Do   FUN[REQ, RES]
}

func NewTask[REQ, RES any, FUN TaskFunc[REQ, RES]](taskFun FUN) *Task[REQ, RES, FUN] {
	return &Task[REQ, RES, FUN]{
		Do: taskFun,
	}
}

type Engine[REQ, RES, V any, FUN TaskFunc[REQ, RES]] struct {
	ctrlEngine   *conctrl.Engine
	Properties   V
	resultChan   chan RES
	resultHandle func(RES)
	failChan     chan *Task[REQ, RES, FUN]
	failHandle   func(*Task[REQ, RES, FUN])
}

func NewEngine[REQ, RES, V any, FUN TaskFunc[REQ, RES]](workerCount uint, props V, fun func(RES), failFun func(*Task[REQ, RES, FUN])) *Engine[REQ, RES, V, FUN] {
	return &Engine[REQ, RES, V, FUN]{
		ctrlEngine:   conctrl.NewEngine(workerCount),
		resultChan:   make(chan RES),
		resultHandle: fun,
		failHandle:   failFun,
		failChan:     make(chan *Task[REQ, RES, FUN]),
		Properties:   props,
	}
}

func (e *Engine[REQ, RES, V, FUN]) NewTaskA(task TaskFuncA) *conctrl.Task {
	return &conctrl.Task{
		Do: task,
	}
}
