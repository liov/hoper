package httpi

import (
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type Body map[string]interface{}

type ResData[T any] struct {
	Code    errorcode.ErrCode `json:"code"`
	Message string            `json:"message,omitempty"`
	//验证码
	Details T `json:"details,omitempty"`
}

func (res *ResData[T]) Response(w http.ResponseWriter, httpcode int) {
	w.WriteHeader(httpcode)
	render.WriteJSON(w, res)
}

func NewResData[T any](code errorcode.ErrCode, msg string, data T) *ResData[T] {
	return &ResData[T]{
		Code:    code,
		Message: msg,
		Details: data,
	}
}

func RespErrcode(w http.ResponseWriter, code errorcode.ErrCode) {
	NewResData(code, code.Error(), nil).Response(w, http.StatusOK)
}

func RespErr(w http.ResponseWriter, err error) {
	NewResData(errorcode.Unknown, err.Error(), nil).Response(w, http.StatusOK)
}

func RespErrMsg(w http.ResponseWriter, msg string) {
	NewResData(errorcode.SUCCESS, msg, nil).Response(w, http.StatusOK)
}

func RespErrRep(w http.ResponseWriter, rep *errorcode.ErrRep) {
	NewResData(rep.Code, rep.Message, nil).Response(w, http.StatusOK)
}

func Response[T any](w http.ResponseWriter, code errorcode.ErrCode, msg string, data T) {
	NewResData(code, msg, data).Response(w, http.StatusOK)
}
