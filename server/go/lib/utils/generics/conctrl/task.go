package conctrl

import (
	"context"
	"time"
)

type Kind uint8

const (
	KindNormal = iota
)

type BaseTaskFunc func(context.Context)

type BaseTask[KEY comparable, T any] struct {
	BaseTaskMeta[KEY]
	BaseTaskFunc
	Props T
}

type BaseTaskMeta[KEY comparable] struct {
	Id       uint64
	Key      KEY
	Priority int
}

func (t *BaseTaskMeta[KEY]) CompareField() int {
	return t.Priority
}

// TODO
type TaskMeta[KEY comparable] struct {
	BaseTaskMeta[KEY]
	Kind Kind
	TaskStatistics
}

func (r *TaskMeta[KEY]) SetKind(k Kind) {
	r.Kind = k
}

func (r *TaskMeta[KEY]) SetKey(key KEY) {
	r.Key = key
}

func (r *TaskMeta[KEY]) SetId(id uint64) {
	r.Id = id
}

type TaskStatistics struct {
	timeCost  time.Duration
	reDoTimes uint
	ErrTimes  int
	Errs      []error
}

type Task[T comparable, P any] struct {
	TaskMeta[T]
	TaskFunc[T, P]
	Props P
}

func (t *Task[KEY, T]) BaseTask() *BaseTask[KEY, T] {
	return &BaseTask[KEY, T]{
		BaseTaskMeta: t.BaseTaskMeta,
		BaseTaskFunc: func(ctx context.Context) {
			t.TaskFunc(ctx)
		},
		Props: t.Props,
	}
}

type Tasks[T comparable, P any] []*Task[T, P]

// ---------------

type ErrHandle func(context.Context, error)

type TaskFunc[T comparable, P any] func(ctx context.Context) ([]*Task[T, P], error)
