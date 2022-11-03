package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type FetchFunc func(ctx context.Context, url string) ([]byte, error)
type ParseFunc func(ctx context.Context, content []byte) ([]conctrl.TaskInterface, error)

type HandleFunc func(ctx context.Context, url string) ([]conctrl.TaskInterface, error)

type Request = conctrl.Task

func NewRequest(key string, taskFunc conctrl.TaskFunc) *Request {
	if taskFunc == nil {
		return nil
	}
	return &Request{TaskMeta: conctrl.TaskMeta{Key: key}, TaskFunc: taskFunc}
}

func NewKindRequest(key string, kind conctrl.Kind, taskFunc conctrl.TaskFunc) *Request {
	if taskFunc == nil {
		return nil
	}
	req := NewRequest(key, taskFunc)
	req.SetKind(kind)
	return req
}

func NewUrlRequest(url string, handleFunc HandleFunc) *Request {
	if handleFunc == nil {
		return nil
	}
	return &Request{TaskMeta: conctrl.TaskMeta{Key: url}, TaskFunc: func(ctx context.Context) ([]conctrl.TaskInterface, error) {
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

func NewHandleFun(f FetchFunc, p ParseFunc) HandleFunc {
	return func(ctx context.Context, url string) ([]conctrl.TaskInterface, error) {
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

type Requests struct {
	reqs       []*Request
	generation int
}

type Engine = conctrl.Engine
