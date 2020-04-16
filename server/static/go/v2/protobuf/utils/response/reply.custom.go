package response

import (
	"encoding/json"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/liov/hoper/go/v2/utils/encoding/protobuf/any"
)

func (m *RawReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return m.Details, nil
}

type GoReply struct {
	Code    uint32
	Message string
	Details proto.Message
}

func (m *AnyReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	var err error
	reply := GoReply{Code: m.Code, Message: m.Message}
	reply.Details, err = any.ResolveAny(m.Details.TypeUrl)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(m.Details.Value, reply.Details)
	if err != nil {
		return nil, err
	}
	return json.Marshal(reply)
}

func (m *HttpResponse) GetContentType() string {
	return "text/template"
}

func (m *HttpResponse) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return m.Body, nil
}

func (m *HttpResponse) Response(w http.ResponseWriter) {
	w.Write(m.Body)
	for k,v:=range m.Header{
		w.Header().Set(k,v)
	}
	w.WriteHeader(int(m.StatusCode))
}
