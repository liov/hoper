package errorcode

import (
	errorsi "github.com/actliboy/hoper/server/go/lib/utils/errors"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

func (x *ErrRep) Error() string {
	return x.Message
}

func (x *ErrRep) GRPCStatus() *status.Status {
	return status.New(codes.Code(x.Code), x.Message)
}

func (x ErrCode) ErrRep() *ErrRep {
	return &ErrRep{Code: x, Message: x.String()}
}

func (x ErrCode) Rep() *ErrRep {
	return &ErrRep{Code: x, Message: x.String()}
}

//example 实现
func (x ErrCode) GRPCStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x ErrCode) Message(msg string) *ErrRep {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) Warp(err error) *ErrRep {
	return &ErrRep{Code: x, Message: err.Error()}
}

func (x ErrCode) Error() string {
	return x.String()
}

func ErrHandle(err interface{}) error {
	if e, ok := err.(*ErrRep); ok {
		return e
	}
	if e, ok := err.(ErrCode); ok {
		return e.ErrRep()
	}
	if e, ok := err.(*status.Status); ok {
		return e.Err()
	}
	if e, ok := err.(*errorsi.ErrRep); ok {
		return e
	}
	if e, ok := err.(errorsi.ErrCode); ok {
		return e.ErrRep()
	}
	if e, ok := err.(error); ok {
		return Unknown.Message(e.Error())
	}
	return Unknown.ErrRep()
}

func Code(err error) int {
	switch v := err.(type) {
	case *ErrRep:
		return int(v.Code)
	case ErrCode:
		return int(v)
	case *errorsi.ErrRep:
		return int(v.Code)
	case errorsi.ErrCode:
		return int(v)
	}
	return 0
}

func (x ErrCode) Origin() errorsi.ErrCode {
	return errorsi.ErrCode(x)
}

func (x ErrCode) OriErrRep() *errorsi.ErrRep {
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: x.String()}
}

func (x ErrCode) OriMessage(msg string) *errorsi.ErrRep {
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: msg}
}

func (x ErrCode) OriWarp(err error) *errorsi.ErrRep {
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: err.Error()}
}

func (x ErrCode) OriLog(err error) *errorsi.ErrRep {
	log.Error(err)
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: x.String()}
}

func (x *ErrRep) MarshalJSON() ([]byte, error) {
	return stringsi.ToBytes(`{"code":` + strconv.Itoa(int(x.Code)) + `,"message":"` + x.Message + `"}`), nil
}

/*func (x ErrCode) MarshalJSON() ([]byte, error) {
	return stringsi.ToBytes(`{"code":` + strconv.Itoa(int(x)) + `,"message":"` + x.String() + `"}`), nil
}

*/

func FromError(err error) (s *ErrRep, ok bool) {
	if err == nil {
		return nil, true
	}
	if se, ok := err.(ErrReInterface); ok {
		return se.ErrRep(), true
	}
	return NewErrReP(codes.Unknown, err.Error()), false
}

func Success() *ErrRep {
	return &ErrRep{Code: SUCCESS, Message: SUCCESS.Error()}
}

func NewErrReP[E ErrCodeGeneric](code E, msg string) *ErrRep {
	return &ErrRep{Code: ErrCode(code), Message: msg}
}

type ErrReInterface interface {
	ErrRep() *ErrRep
}

type ErrCodeInterface interface {
}

type ErrCodeGeneric interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}
