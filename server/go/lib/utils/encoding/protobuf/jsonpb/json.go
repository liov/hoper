package jsonpb

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/liov/hoper/server/go/lib/utils/encoding/json/iterator"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
	"github.com/liov/hoper/server/go/lib/protobuf/response"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
)

var JsonPb = JSONPb{iterator.Standard}

type JSONPb struct {
	jsoniter.API
}

func (*JSONPb) ContentType(_ interface{}) string {
	return "application/json"
}

func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	if _, ok := v.(error); ok {
		return j.API.Marshal(v)
	}
	if msg, ok := v.(*wrappers.StringValue); ok {
		v = msg.Value
	}
	return j.API.Marshal(&httpi.ResData{
		Code:    0,
		Message: "OK",
		Details: v,
	})
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
