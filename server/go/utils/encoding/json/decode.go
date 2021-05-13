package json

import (
	"encoding/json"
	"io"
)

func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}
