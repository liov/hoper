package errorsi

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"github.com/gin-gonic/gin/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"sync"
)

var codeMap = sync.Map{}

type DefaultErrRep interface {
	ErrRep() *ErrRep
}

type GRPCStatus interface {
	GRPCStatus() *status.Status
}

type ErrCode uint32

func SetCode(code ErrCode, msg string) {
	codeMap.Store(code, msg)
}

func (x ErrCode) String() string {
	value, ok := codeMap.Load(x)
	if ok {
		return value.(string)
	}
	return ""
}

func (x ErrCode) ErrRep() *ErrRep {
	return &ErrRep{Code: x, Message: x.String()}
}

// example 实现
func (x ErrCode) GRPCStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x ErrCode) Message(msg string) *ErrRep {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) MarshalJSON() ([]byte, error) {
	return stringsi.ToBytes(`{"code":` + strconv.Itoa(int(x)) + `,"message":"` + x.String() + `"}`), nil
}

func (x ErrCode) Warp(err error) *ErrRep {
	return &ErrRep{Code: x, Message: err.Error()}
}

func (x ErrCode) Log(err error) *ErrRep {
	log.Error(err)
	return &ErrRep{Code: x, Message: x.String()}
}

func (x ErrCode) Error() string {
	return x.String()
}

func (x ErrCode) Response(w http.ResponseWriter) {
	render.WriteJSON(w, x)
}
