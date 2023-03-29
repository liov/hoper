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
	id       uint64
	Key      KEY
	Describe string
	Priority int
}

func (t *BaseTaskMeta[KEY]) CmpKey() int {
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

func (r *TaskMeta[KEY]) Id() uint64 {
	return r.id
}

type TaskStatistics struct {
	timeCost  time.Duration
	reDoTimes uint
	ErrTimes  int
}

type Task[KEY comparable, P any] struct {
	TaskMeta[KEY]
	TaskFunc[KEY, P]
	errs  []error
	Props P
}

func (t *Task[KEY, P]) Errs() []error {
	return t.errs
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
