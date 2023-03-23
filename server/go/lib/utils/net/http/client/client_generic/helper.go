package client_generic

import "github.com/liov/hoper/server/go/lib/utils/net/http/client"

func SimpleGet[RES any](url string) (*RES, error) {
	return NewGetRequest[RES](url).Do(nil)
}

func GetSubData[RES GetDataInterface[T], T any](url string) (*T, error) {
	var response RES
	err := new(client.Request).Get(url, &response)
	if err != nil {
		return new(T), err
	}
	return response.GetData(), nil
}

// Deprecated
func CustomizeGet[RES GetDataInterface[T], T any](setParams client.SetParams) func(url string) (*T, error) {
	return func(url string) (*T, error) {
		var response RES
		err := setParams(new(client.Request)).Get(url, &response)
		if err != nil {
			return new(T), err
		}
		return response.GetData(), nil
	}
}
