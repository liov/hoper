package client_generic

import (
	"errors"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"

	"net/http"
	"time"
)

// RequestWrapper ...
type RequestWrapper[RES any] client.Request

func NewRequest[RES any](url, method string) *RequestWrapper[RES] {
	return (*RequestWrapper[RES])(client.NewRequest(url, method))
}

func (req *RequestWrapper[RES]) ContentType(contentType client.ContentType) *RequestWrapper[RES] {
	(*client.Request)(req).ContentType(contentType)
	return req
}

func (req *RequestWrapper[RES]) SetHeader(header client.Header) *RequestWrapper[RES] {
	(*client.Request)(req).SetHeader(header)
	return req
}

func (req *RequestWrapper[RES]) AddHeader(k, v string) *RequestWrapper[RES] {
	(*client.Request)(req).AddHeader(k, v)
	return req
}

func (req *RequestWrapper[RES]) SetLogger(logger client.LogCallback) *RequestWrapper[RES] {
	(*client.Request)(req).WithLogger(logger)
	return req
}

func (req *RequestWrapper[RES]) ResponseHandler(handler func([]byte) ([]byte, error)) *RequestWrapper[RES] {
	(*client.Request)(req).ResponseHandler(handler)
	return req
}

func (req *RequestWrapper[RES]) Timeout(timeout time.Duration) *RequestWrapper[RES] {
	(*client.Request)(req).Timeout(timeout)
	return req
}

func (req *RequestWrapper[RES]) WithClient(c *http.Client) *RequestWrapper[RES] {
	(*client.Request)(req).WithClient(c)
	return req
}

func (req *RequestWrapper[RES]) RetryTimes(times int) *RequestWrapper[RES] {
	(*client.Request)(req).RetryTimes(times)
	return req
}
func (req *RequestWrapper[RES]) DisableLog() *RequestWrapper[RES] {
	(*client.Request)(req).DisableLog()
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
func (r *RequestWrapper[RES]) Do(req any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Do(req, response)
	return response, err
}

func (r *RequestWrapper[RES]) Get(url string) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Get(url, response)
	return response, err
}

func (r *RequestWrapper[RES]) Post(url string, param any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Post(url, param, response)
	return response, err
}

func NewGetRequest[RES any](url string) *RequestWrapper[RES] {
	return NewRequest[RES](url, http.MethodGet)
}

type SetParams[RES any] func(req *RequestWrapper[RES]) *RequestWrapper[RES]
