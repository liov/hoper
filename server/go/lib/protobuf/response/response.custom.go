package response

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

type GoReply struct {
	Code    uint32
	Message string
	Details proto.Message
}

func (x *HttpResponse) GetContentType() string {
	return x.Header["Content-Type"]
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

var ResponseOk = &TinyRep{Message: "OK"}
