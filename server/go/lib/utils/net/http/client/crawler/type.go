package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type FetchFun func(ctx context.Context, url string) ([]byte, error)
type ParseFun func(ctx context.Context, content []byte) ([]*Request, error)

type HandleFun func(ctx context.Context, url string) ([]*Request, error)
type TaskFun func(context.Context) ([]*Request, error)

type Request struct {
	conctrl.TaskMeta
	Key      string
	TaskFun  TaskFun
	errTimes int
}

func NewRequest(key string, taskFun TaskFun) *Request {
	if taskFun == nil {
		return nil
	}
	return &Request{Key: key, TaskFun: taskFun}
}

func NewKindRequest(key string, kind conctrl.Kind, taskFun TaskFun) *Request {
	if taskFun == nil {
		return nil
	}
	return NewRequest(key, taskFun).SetKind(kind)
}

func NewUrlRequest(url string, handleFun HandleFun) *Request {
	if handleFun == nil {
		return nil
	}
	return &Request{Key: url, TaskFun: func(ctx context.Context) ([]*Request, error) {
		return handleFun(ctx, url)
	}}
}

func NewUrlKindRequest(url string, kind conctrl.Kind, handleFun HandleFun) *Request {
	if handleFun == nil {
		return nil
	}
	return NewUrlRequest(url, handleFun).SetKind(kind)
}

func (r *Request) SetKind(k conctrl.Kind) *Request {
	r.Kind = k
	return r
}

func (r *Request) SetKey(key string) *Request {
	r.Key = key
	return r
}

func (r *Request) SetId(id uint) *Request {
	r.Id = id
	return r
}

func NewHandleFun(f FetchFun, p ParseFun) HandleFun {
	return func(ctx context.Context, url string) ([]*Request, error) {
		content, err := f(ctx, url)
		if err != nil {
			return nil, err
		}
		return p(ctx, content)
	}
}

func NewUrlRequest2(url string, fetchFun FetchFun, parseFunction ParseFun) *Request {
	return NewUrlRequest(url, NewHandleFun(fetchFun, parseFunction))
}

type ReqInterface interface {
	HandleFun
}

type HandleFuncs []HandleFun

func (h HandleFun) Append(handleFun HandleFun) *HandleFuncs {
	newh := append(HandleFuncs{h}, handleFun)
	return &newh
}

func (h *HandleFuncs) Append(handleFun HandleFun) *HandleFuncs {
	newh := append(*h, handleFun)
	return &newh
}
