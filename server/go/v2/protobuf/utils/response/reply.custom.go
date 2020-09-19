package response

import (
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

func (m *RawReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return *m.Details, nil
}

type GoReply struct {
	Code    uint32
	Message string
	Details proto.Message
}

func (m *AnyReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	var err error
	reply := GoReply{Code: m.Code, Message: m.Message}
	reply.Details, err = m.Details.UnmarshalNew()
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
	//我尼玛也是头一次知道要按顺序来的 response.wroteHeader
	//先设置请求头，再设置状态码，再写body
	//原因是http里每次操作都要判断wroteHeader(表示已经写过header了，不可以再写了)
	for k, v := range m.Header {
		w.Header().Set(k, v)
	}
	w.WriteHeader(int(m.StatusCode))
	w.Write(m.Body)
}
