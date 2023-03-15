package httpi

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	"io"
	"net/http"
	"time"
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

type ResAnyData = ResData[any]

func NewResData[T any](code errorcode.ErrCode, msg string, data T) *ResData[T] {
	return &ResData[T]{
		Code:    code,
		Message: msg,
		Details: data,
	}
}

func RespErrcode(w http.ResponseWriter, code errorcode.ErrCode) {
	NewResData[any](code, code.Error(), nil).Response(w, http.StatusOK)
}

func RespErr(w http.ResponseWriter, err error) {
	NewResData[any](errorcode.Unknown, err.Error(), nil).Response(w, http.StatusOK)
}

func RespErrMsg(w http.ResponseWriter, msg string) {
	NewResData[any](errorcode.SUCCESS, msg, nil).Response(w, http.StatusOK)
}

func RespErrRep(w http.ResponseWriter, rep *errorcode.ErrRep) {
	NewResData[any](rep.Code, rep.Message, nil).Response(w, http.StatusOK)
}

func Response[T any](w http.ResponseWriter, code errorcode.ErrCode, msg string, data T) {
	NewResData(code, msg, data).Response(w, http.StatusOK)
}

func StreamWriter(w http.ResponseWriter, writer func(w io.Writer) bool) {
	notifyClosed := w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		// response writer forced to close, exit.
		case <-notifyClosed:
			return
		default:
			shouldContinue := writer(w)
			w.(http.Flusher).Flush()
			if !shouldContinue {
				return
			}
		}
	}
}

func Stream(w http.ResponseWriter) {
	w.Header().Set("X-Accel-Buffering", "no") //nginx的锅必须加
	w.Header().Set("Transfer-Encoding", "chunked")
	i := 0
	ints := []int{1, 2, 3, 5, 7, 9, 11, 13, 15, 17, 23, 29}
	StreamWriter(w, func(w io.Writer) bool {
		fmt.Fprintf(w, "Message number %d<br>", ints[i])
		time.Sleep(500 * time.Millisecond) // simulate delay.
		if i == len(ints)-1 {
			return false //关闭并刷新
		}
		i++
		return true //继续写入数据
	})
}

var ResponseSysErr = []byte(`{"code":10000,"message":"系统错误"}`)
var ResponseOk = []byte(`{"code":0,"message":"OK"}`)
