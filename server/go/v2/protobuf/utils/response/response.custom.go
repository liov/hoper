package response

import (
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

func (x *RawReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return *x.Details, nil
}

type GoReply struct {
	Code    uint32
	Message string
	Details proto.Message
}

func (x *AnyReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	var err error
	reply := GoReply{Code: x.Code, Message: x.Message}
	reply.Details, err = x.Details.UnmarshalNew()
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(x.Details.Value, reply.Details)
	if err != nil {
		return nil, err
	}
	return json.Marshal(reply)
}

func (x *HttpResponse) GetContentType() string {
	return x.Header["Content-Type"]
}

func (x *HttpResponse) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return x.Body, nil
}

func (x *HttpResponse) Response(w http.ResponseWriter) {
	//我尼玛也是头一次知道要按顺序来的 response.wroteHeader
	//先设置请求头，再设置状态码，再写body
	//原因是http里每次操作都要判断wroteHeader(表示已经写过header了，不可以再写了)
	for k, v := range x.Header {
		w.Header().Set(k, v)
	}
	w.WriteHeader(int(x.StatusCode))
	w.Write(x.Body)
}
