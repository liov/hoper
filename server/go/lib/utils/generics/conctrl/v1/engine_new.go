package v1

import (
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

func NewTask[REQ, RES any, FUN TaskFunc[REQ, RES]](taskFun FUN) *Task[REQ, RES, FUN] {
	return &Task[REQ, RES, FUN]{
		Do: taskFun,
	}
}

type Engine[REQ, RES, V any, FUN TaskFunc[REQ, RES]] struct {
	ctrlEngine   *conctrl.BaseEngine
	Properties   V
	resultChan   chan RES
	resultHandle func(RES)
	failChan     chan *Task[REQ, RES, FUN]
	failHandle   func(*Task[REQ, RES, FUN])
}

func NewEngineN[REQ, RES, V any, FUN TaskFunc[REQ, RES]](workerCount uint, props V, fun func(RES), failFun func(*Task[REQ, RES, FUN])) *Engine[REQ, RES, V, FUN] {
	return &Engine[REQ, RES, V, FUN]{
		ctrlEngine:   conctrl.NewBaseEngine(workerCount),
		resultChan:   make(chan RES),
		resultHandle: fun,
		failHandle:   failFun,
		failChan:     make(chan *Task[REQ, RES, FUN]),
		Properties:   props,
	}
}

func (e *Engine[REQ, RES, V, FUN]) NewTaskA(task TaskFuncA) *conctrl.Task {
	return &conctrl.Task{
		Do: conctrl.BaseTaskFunc(task),
	}
}
