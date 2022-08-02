package crawler

import "github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"

type FetchFun[T any] func(url string) (T, error)
type ParseFun[T any] func(t T) ([]*crawler.Request, error)

func NewRequest(url string, handleFun crawler.HandleFun) *crawler.Request {
	return &crawler.Request{Url: url, HandleFun: handleFun}
}

func NewHandleFun[T any](f FetchFun[T], p ParseFun[T]) crawler.HandleFun {
	return func(url string) ([]*crawler.Request, error) {
		content, err := f(url)
		if err != nil {
			return nil, err
		}
		return p(content)
	}
}

func NewRequest2[T any](url string, fetchFun FetchFun[T], parseFunction ParseFun[T]) *crawler.Request {
	return &crawler.Request{Url: url, HandleFun: NewHandleFun[T](fetchFun, parseFunction)}
}

type Callback[T any] func(t T) error
