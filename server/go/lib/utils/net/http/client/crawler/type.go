package crawler

type FetchFun func(url string) ([]byte, error)
type ParseFun func(content []byte) ([]*Request, error)
type FailFunc func(err error)

type HandleFun func(url string) ([]*Request, error)

type Request struct {
	Url       string
	HandleFun HandleFun
	FailFun   FailFunc
}

func NewRequest(url string, handleFun HandleFun, parseFunction ParseFun) *Request {
	return &Request{Url: url, HandleFun: handleFun}
}

func NewHandleFun(f FetchFun, p ParseFun) HandleFun {
	return func(url string) ([]*Request, error) {
		content, err := f(url)
		if err != nil {
			return nil, err
		}
		return p(content)
	}
}

func NewRequest2(url string, fetchFun FetchFun, parseFunction ParseFun) *Request {
	return &Request{Url: url, HandleFun: NewHandleFun(fetchFun, parseFunction)}
}

type ReqInterface interface {
	HandleFun
}
