package client

import (
	"io"
	"net/http"
)

func SimpleGet(url string, response any) error {
	return NewGetRequest(url).DoWithNoParam(response)
}

func SimpleGetStream(url string) (io.ReadCloser, error) {
	var resp *http.Response
	err := SimpleGet(url, &resp)
	if err != nil {
		return resp.Body, err
	}
	return resp.Body, nil
}

func SimplePost(url string, param, response interface{}) error {
	return NewPostRequest(url).Do(param, response)
}

func SimplePut(url string, param, response interface{}) error {
	return NewPutRequest(url).Do(param, response)
}

func SimpleDelete(url string, param, response interface{}) error {
	return NewDeleteRequest(url).Do(param, response)
}
