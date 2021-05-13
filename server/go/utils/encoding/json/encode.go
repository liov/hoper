package json

import (
	"encoding/json"
	"io"
)

func NewEncoder(r io.Writer) *json.Encoder {
	return json.NewEncoder(r)
}
