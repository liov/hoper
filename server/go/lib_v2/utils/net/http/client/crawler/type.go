package crawler

import (
	"context"
	"github.com/liov/hoper/server/go/lib_v2/utils/conctrl"
)

type Prop struct {
}

type Request = conctrl.Task[string, Prop]
type BaseTaskMeta = conctrl.BaseTaskMeta[string]
type TaskMeta = conctrl.TaskMeta[string]
type TaskFunc = conctrl.TaskFunc[string, Prop]

func NewRequest(key string, kind conctrl.Kind, taskFunc TaskFunc) *Request {
	return &Request{
		TaskMeta: TaskMeta{
			BaseTaskMeta: BaseTaskMeta{Key: key},
			Kind:         kind,
		},
		TaskFunc: taskFunc,
	}
}

type Engine = conctrl.Engine[string, Prop, Prop]

func NewEngine(workerCount uint) *conctrl.Engine[string, Prop, Prop] {
	return conctrl.NewEngine[string, Prop, Prop](workerCount)
}

type HandleFunc func(ctx context.Context, url string) ([]*Request, error)

func NewUrlRequest(url string, handleFunc HandleFunc) *Request {
	if handleFunc == nil {
		return nil
	}
	return &Request{TaskMeta: TaskMeta{BaseTaskMeta: BaseTaskMeta{Key: url}}, TaskFunc: func(ctx context.Context) ([]*Request, error) {
		return handleFunc(ctx, url)
	}}
}

func NewUrlKindRequest(url string, kind conctrl.Kind, handleFunc HandleFunc) *Request {
	if handleFunc == nil {
		return nil
	}
	req := NewUrlRequest(url, handleFunc)
	req.SetKind(kind)
	return req
}

func NewTaskMeta(key string) TaskMeta {
	return TaskMeta{BaseTaskMeta: BaseTaskMeta{Key: key}}
}
