package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
)

type FetchFun[T any] func(ctx context.Context, url string) (T, error)
type ParseFun[T any] func(ctx context.Context, t T) ([]*crawler.Request, error)

func NewRequest(url string, handleFun crawler.HandleFun) *crawler.Request {
	return &crawler.Request{Url: url, HandleFun: handleFun}
}

func NewHandleFun[T any](f FetchFun[T], p ParseFun[T]) crawler.HandleFun {
	return func(ctx context.Context, url string) ([]*crawler.Request, error) {
		content, err := f(ctx, url)
		if err != nil {
			return nil, err
		}
		return p(ctx, content)
	}
}

func NewRequest2[T any](url string, fetchFun FetchFun[T], parseFunction ParseFun[T]) *crawler.Request {
	return &crawler.Request{Url: url, HandleFun: NewHandleFun[T](fetchFun, parseFunction)}
}

type Callback[T any] func(t T) error