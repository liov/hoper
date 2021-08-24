package httpi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/render"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
)

type H map[string]interface{}

type ResData struct {
	Code    uint32 `json:"code"`
	Message string `json:"message,omitempty"`
	//验证码
	Details interface{} `json:"details,omitempty"`
}

func (res *ResData) Response(w http.ResponseWriter) {
	render.WriteJSON(w, res)
}

//先信息后数据最后状态码
//入参1. data interface{},msg string,code int
//2.msg,code |data默认nil
//3.data,msg |code默认SUCCESS
//4.msg |data默认nil code默认ERROR
//5.data |msg默认"",code默认SUCCESS
func Response(w http.ResponseWriter, res ...interface{}) {

	var resData ResData

	if len(res) == 1 {
		resData.Code = uint32(errorcode.Unknown)
		if msgTmp, ok := res[0].(string); ok {
			resData.Message = msgTmp
			resData.Details = nil
		} else {
			resData.Details = res[0]
			resData.Code = uint32(errorcode.SUCCESS)
		}
	} else if len(res) == 2 {
		if msgTmp, ok := res[0].(string); ok {
			resData.Details = nil
			resData.Message = msgTmp
			resData.Code = res[1].(uint32)
		} else {
			resData.Details = res[0]
			resData.Message = res[1].(string)
			resData.Code = uint32(errorcode.SUCCESS)
		}
	} else {
		resData.Details = res[0]
		resData.Message = res[1].(string)
		resData.Code = res[2].(uint32)
	}

	render.WriteJSON(w, &resData)
}

func Res(w http.ResponseWriter, code uint32, msg string, data interface{}) {
	var resData = ResData{
		Code:    code,
		Message: msg,
		Details: data,
	}
	render.WriteJSON(w, &resData)
}

type File struct {
	File http.File
	Name string
}

type HttpFile interface {
	io.Reader
	Name() string
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
