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

type BaseTask[T any] struct {
	BaseTaskMeta
	BaseTaskFunc
	Props T
}

type BaseTaskMeta struct {
	Id       uint64
	Priority int
}

func (t *BaseTaskMeta) CompareField() int {
	return t.Priority
}

// TODO
type TaskMeta[T comparable] struct {
	BaseTaskMeta
	Kind Kind
	Key  T
	TaskStatistics
}

func (r *TaskMeta[T]) SetKind(k Kind) {
	r.Kind = k
}

func (r *TaskMeta[T]) SetKey(key T) {
	r.Key = key
}

func (r *TaskMeta[T]) SetId(id uint64) {
	r.Id = id
}

type TaskStatistics struct {
	timeCost  time.Duration
	reDoTimes uint
	ErrTimes  int
}

type Task[T comparable, P any] struct {
	TaskMeta[T]
	TaskFunc[T, P]
	Props P
}

func (t *Task[T, P]) BaseTask() *BaseTask[P] {
	return &BaseTask[P]{
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
