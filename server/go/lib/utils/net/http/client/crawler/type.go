package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type FetchFun func(ctx context.Context, url string) ([]byte, error)
type ParseFun func(ctx context.Context, content []byte) ([]*Request, error)

type HandleFun func(ctx context.Context, url string) ([]*Request, error)

type Request struct {
	Id        uint
	Kind      conctrl.Kind
	Url       string
	HandleFun HandleFun
}

func NewRequest(url string, handleFun HandleFun) *Request {
	return &Request{Url: url, HandleFun: handleFun}
}

func NewKindRequest(url string, kind conctrl.Kind, handleFun HandleFun) *Request {
	return &Request{Url: url, Kind: kind, HandleFun: handleFun}
}

func (r *Request) SetKind(k conctrl.Kind) *Request {
	r.Kind = k
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

func NewRequest2(url string, fetchFun FetchFun, parseFunction ParseFun) *Request {
	return &Request{Url: url, HandleFun: NewHandleFun(fetchFun, parseFunction)}
}

type ReqInterface interface {
	HandleFun
}

func (r *Request) NewTask() *conctrl.Task {
	return &conctrl.Task{Id: r.Id, Kind: r.Kind, Do: func(ctx context.Context) {
		r.HandleFun(ctx, r.Url)
	}}
}
