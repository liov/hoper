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
