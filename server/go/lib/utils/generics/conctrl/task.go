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

func (t *BaseTaskMeta[KEY]) SetPriority(priority int) {
	t.Priority = priority
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
}

type Task[KEY comparable, P any] struct {
	TaskMeta[KEY]
	TaskFunc[KEY, P]
	Errs  []error
	Props P
}

func (t *Task[KEY, P]) BaseTask(handle func(tasks []*Task[KEY, P], err error)) *BaseTask[KEY, P] {
	return &BaseTask[KEY, P]{
		BaseTaskMeta: t.BaseTaskMeta,
		BaseTaskFunc: func(ctx context.Context) {
			handle(t.TaskFunc(ctx))
		},
		Props: t.Props,
	}
}

type Tasks[KEY comparable, P any] []*Task[KEY, P]

// ---------------

type ErrHandle func(context.Context, error)

type TaskFunc[KEY comparable, P any] func(ctx context.Context) ([]*Task[KEY, P], error)
