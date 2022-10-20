package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type FetchFunc func(ctx context.Context, url string) ([]byte, error)
type ParseFunc func(ctx context.Context, content []byte) ([]*Request, error)

type HandleFunc func(ctx context.Context, url string) ([]*Request, error)
type TaskFunc func(context.Context) ([]*Request, error)

type TaskFuncInterface interface {
	TaskFunc(context.Context) ([]*Request, error)
}

func (t TaskFunc) TaskFunc(ctx context.Context) ([]*Request, error) {
	return t(ctx)
}

type Request struct {
	conctrl.TaskMeta
	Key      string
	TaskFunc TaskFunc
	errTimes int
}

func NewRequest(key string, taskFunc TaskFunc) *Request {
	if taskFunc == nil {
		return nil
	}
	return &Request{Key: key, TaskFunc: taskFunc}
}

func NewKindRequest(key string, kind conctrl.Kind, taskFunc TaskFunc) *Request {
	if taskFunc == nil {
		return nil
	}
	return NewRequest(key, taskFunc).SetKind(kind)
}

func NewUrlRequest(url string, handleFunc HandleFunc) *Request {
	if handleFunc == nil {
		return nil
	}
	return &Request{Key: url, TaskFunc: func(ctx context.Context) ([]*Request, error) {
		return handleFunc(ctx, url)
	}}
}

func NewUrlKindRequest(url string, kind conctrl.Kind, handleFunc HandleFunc) *Request {
	if handleFunc == nil {
		return nil
	}
	return NewUrlRequest(url, handleFunc).SetKind(kind)
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

func NewHandleFun(f FetchFunc, p ParseFunc) HandleFunc {
	return func(ctx context.Context, url string) ([]*Request, error) {
		content, err := f(ctx, url)
		if err != nil {
			return nil, err
		}
		return p(ctx, content)
	}
}

func NewUrlRequest2(url string, fetchFunc FetchFunc, parseFunc ParseFunc) *Request {
	return NewUrlRequest(url, NewHandleFun(fetchFunc, parseFunc))
}

type RequestInterface interface {
	HandleFunc
}

type HandleFuncs []HandleFunc

func (h HandleFunc) Append(handleFunc HandleFunc) *HandleFuncs {
	newh := append(HandleFuncs{h}, handleFunc)
	return &newh
}

func (h *HandleFuncs) Append(handleFunc HandleFunc) *HandleFuncs {
	newh := append(*h, handleFunc)
	return &newh
}
