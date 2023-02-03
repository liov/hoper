package conctrl

import (
	"context"
	"time"
)

type Kind = uint8

const (
	KindNormal = iota
)

type BaseTaskFunc func(context.Context)

type BaseTask struct {
	BaseTaskMeta
	BaseTaskFunc
}

type BaseTaskHeap []*BaseTask

func (b *BaseTaskHeap) Less(i, j int) bool {
	return (*b)[i].Priority < (*b)[j].Priority
}

func (b *BaseTaskHeap) Swap(i, j int) {
	(*b)[i], (*b)[j] = (*b)[j], (*b)[i]
}

func (b *BaseTaskHeap) Push(x any) {
	*b = append(*b, x.(*BaseTask))
}

func (b *BaseTaskHeap) Pop() any {
	v := (*b)[b.Len()-1]
	*b = (*b)[:b.Len()-1]
	return v
}

func (b *BaseTaskHeap) Len() int {
	return len(*b)
}

type BaseTaskMeta struct {
	Id       uint64
	Key      string
	Priority int
}

func (t *BaseTaskMeta) CompareField() int {
	return t.Priority
}

// TODO
type TaskMeta struct {
	BaseTaskMeta
	Kind Kind
	TaskStatistics
}

func (r *TaskMeta) SetKind(k Kind) {
	r.Kind = k
}

func (r *TaskMeta) SetKey(key string) {
	r.Key = key
}

func (r *TaskMeta) SetId(id uint64) {
	r.Id = id
}

type TaskStatistics struct {
	timeCost  time.Duration
	reDoTimes uint
	ErrTimes  int
}

type TaskFunc1 func(ctx context.Context) ([]*Task, error)

type Task struct {
	TaskMeta
	TaskFunc
	Errs []error
}

func (t *Task) HasTask() *Task {
	return t
}

func (t *Task) BaseTask() *BaseTask {
	return &BaseTask{
		BaseTaskMeta: t.BaseTaskMeta,
		BaseTaskFunc: func(ctx context.Context) {
			t.TaskFunc(ctx)
		},
	}
}

type Tasks []*Task

// ---------------

type ErrHandle func(context.Context, error)

type TaskInterface interface {
	HasTask() *Task
}

type TaskFunc func(ctx context.Context) ([]TaskInterface, error)

type TaskFuncInterface interface {
	TaskFunc(context.Context) ([]TaskInterface, error)
}

func (t TaskFunc) TaskFunc(ctx context.Context) ([]TaskInterface, error) {
	return t(ctx)
}
