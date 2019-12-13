package json

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"
)

type JSONPb struct {
	jsoniter.API
}

func (*JSONPb) ContentType() string {
	return "application/json"
}

func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := j.marshalTo(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *JSONPb) marshalTo(w io.Writer, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := j.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}
	return (*jsonpb.Marshaler)(j).Marshal(w, p)
}
