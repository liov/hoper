package client_generic

import "github.com/liov/hoper/server/go/lib/utils/net/http/client"

type GetDataInterface[T any] interface {
	GetData() *T
}

type SubDataRequest[RES GetDataInterface[T], T any] RequestWrapper[RES]

func NewSubDataRequestParams[RES GetDataInterface[T], T any](url, method string) *SubDataRequest[RES, T] {
	return (*SubDataRequest[RES, T])(NewRequest[RES](url, method))
}

func (req *SubDataRequest[RES, T]) Origin() *client.Request {
	return (*client.Request)(req)
}

func (req *SubDataRequest[RES, T]) OriginWrapper() *RequestWrapper[GetDataInterface[T]] {
	return (*RequestWrapper[GetDataInterface[T]])(req)
}

// Do create a HTTP request
func (r *SubDataRequest[RES, T]) Do(req any) (*RES, error) {
	response := new(RES)
	err := (*client.Request)(r).Do(req, response)
	return response, err
}

func (req *SubDataRequest[RES, T]) Get(url string) (*T, error) {
	var response RES
	err := (*client.Request)(req).Url(url).Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.GetData(), nil
}
