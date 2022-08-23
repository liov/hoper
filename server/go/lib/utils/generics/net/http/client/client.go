package client

import (
	"errors"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"

	"net/http"
	"time"
)

// RequestParams ...
type RequestParams[RES any] client.RequestParams

func NewRequest[RES any](url, method string) *RequestParams[RES] {
	return (*RequestParams[RES])(client.NewRequest(url, method))
}

func (req *RequestParams[RES]) ContentType(contentType client.ContentType) *RequestParams[RES] {
	(*client.RequestParams)(req).ContentType(contentType)
	return req
}

func (req *RequestParams[RES]) AddHeader(k, v string) *RequestParams[RES] {
	(*client.RequestParams)(req).AddHeader(k, v)
	return req
}

func (req *RequestParams[RES]) SetLogger(logger client.LogCallback) *RequestParams[RES] {
	(*client.RequestParams)(req).WithLogger(logger)
	return req
}

func (req *RequestParams[RES]) ResponseHandler(handler func([]byte) ([]byte, error)) *RequestParams[RES] {
	(*client.RequestParams)(req).ResponseHandler(handler)
	return req
}

func (req *RequestParams[RES]) Timeout(timeout time.Duration) *RequestParams[RES] {
	(*client.RequestParams)(req).Timeout(timeout)
	return req
}

func (req *RequestParams[RES]) WithClient(c *http.Client) *RequestParams[RES] {
	(*client.RequestParams)(req).WithClient(c)
	return req
}

type ResponseBody[RES any] struct {
	Status  int    `json:"status"`
	Data    RES    `json:"data"`
	Message string `json:"message"`
}

func CommonResponse[RES any]() client.ResponseBodyCheck {
	return &ResponseBody[RES]{}
}

func (res *ResponseBody[RES]) CheckError() error {
	if res.Status != 0 {
		return errors.New(res.Message)
	}
	return nil
}

// Do create a HTTP request
func (r *RequestParams[RES]) Do(req any) (*RES, error) {
	response := new(RES)
	err := (*client.RequestParams)(r).Do(req, response)
	return response, err
}

func Get[RES any](url string) (*RES, error) {
	return NewGetRequest[RES](url).Do(nil)
}

func NewGetRequest[RES any](url string) *RequestParams[RES] {
	return NewRequest[RES](url, http.MethodGet)
}
