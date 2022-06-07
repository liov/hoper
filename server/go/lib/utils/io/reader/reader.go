package reader

import (
	"io"
	"io/ioutil"
)

func ReadCloser(body io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	body.Close()
	return data, nil
}
