package reader

import (
	"io"
)

func ReadCloser(body io.ReadCloser) ([]byte, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	body.Close()
	return data, nil
}
