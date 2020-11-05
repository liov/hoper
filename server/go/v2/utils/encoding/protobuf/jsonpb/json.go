package jsonpb

import (
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
)

type JSONPb struct {
	jsoniter.API
}

func (*JSONPb) ContentType(_ interface{}) string {
	return "application/json"
}

func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	if rb, ok := v.(jsonpb.JSONPBMarshaler); ok {
		return rb.MarshalJSONPB(nil)
	}
	return j.API.Marshal(v)
}

func (j *JSONPb) Unmarshal(data []byte, v interface{}) error {
	return j.API.Unmarshal(data, v)
}

func (j *JSONPb) Delimiter() []byte {
	return []byte("\n")
}

// NewDecoder returns a runtime.Decoder which reads JSON stream from "r".
func (j *JSONPb) NewDecoder(r io.Reader) runtime.Decoder {
	return j.API.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPb) NewEncoder(w io.Writer) runtime.Encoder {
	return j.API.NewEncoder(w)
}

func (j *JSONPb) ContentTypeFromMessage(v interface{}) string {
	if httpBody, ok := v.(*response.HttpResponse); ok {
		return httpBody.GetContentType()
	}
	return j.ContentType(v)
}
