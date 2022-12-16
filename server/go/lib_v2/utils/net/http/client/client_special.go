package client

type GetDataInterface[T any] interface {
	GetData() T
}

type SubDataRequestParams[RES GetDataInterface[T], T any] RequestParams[RES]

func NewSubDataRequestParams[RES GetDataInterface[T], T any](url, method string) *SubDataRequestParams[RES, T] {
	return (*SubDataRequestParams[RES, T])(NewRequest(url, method))
}

func (req *SubDataRequestParams[RES, T]) Origin() *RequestParams[GetDataInterface[T]] {
	return (*RequestParams[GetDataInterface[T]])(req)
}

func (req *SubDataRequestParams[RES, T]) GetSubData() (T, error) {
	res, err := req.Origin().Do(nil)
	if err != nil {
		return nil, err
	}
	// 这种实现...,RES实现了GetData, &RES不应该判断为自动实现吗
	return (*res).GetData(), nil
}
