package new

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/conctrl"
)

type TaskFuncA func(context.Context)

type TaskFuncB func(context.Context) error

type TaskFuncC[T any] func(context.Context) (T, error)

type TaskFuncD[T any] func(context.Context, T) error

type TaskFuncE[REQ, RES any] func(context.Context, REQ) (RES, error)

type TaskFunc[REQ, RES any] interface {
	TaskFuncA | TaskFuncB | TaskFuncC[RES] | TaskFuncD[REQ] | TaskFuncE[REQ, RES]
}

type Task[REQ, RES any, FUN TaskFunc[REQ, RES]] struct {
	Id   uint
	Kind conctrl.Kind
	Do   FUN[REQ, RES]
}
