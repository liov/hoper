package errorsi

import (
	"github.com/gin-gonic/gin/render"
	"github.com/liov/hoper/go/v2/utils/log"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

type ErrCode uint32

func (x ErrCode) String() string {
	return ""
}

func (x *ErrRep) Error() string {
	return x.Message
}

func (x *ErrRep) GRPCStatus() *status.Status {
	return status.New(codes.Code(x.Code), x.Message)
}

func (x ErrCode) ErrRep() *ErrRep {
	return &ErrRep{Code: x, Message: x.String()}
}

//example 实现
func (x ErrCode) GRPCStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x ErrCode) Message(msg string) *ErrRep {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) Warp(err error)  *ErrRep {
	return &ErrRep{Code: x, Message: err.Error()}
}

func (x ErrCode) Log(err error)  *ErrRep {
	log.Error(err)
	return &ErrRep{Code: x, Message: x.String()}
}

func (x ErrCode) Error() string {
	return x.String()
}

func (x *ErrRep) MarshalJSON() ([]byte,error){
	return stringsi.ToBytes(`{"code":`+strconv.Itoa(int(x.Code))+`,"message":`+x.Message+`}`),nil
}

func (x ErrCode) MarshalJSON() ([]byte,error){
	return stringsi.ToBytes(`{"code":`+strconv.Itoa(int(x))+`,"message":`+x.String()+`}`),nil
}


func (x *ErrRep) Response(w http.ResponseWriter){
	render.WriteJSON(w, x)
}

func (x ErrCode) Response(w http.ResponseWriter){
	render.WriteJSON(w, x)
}