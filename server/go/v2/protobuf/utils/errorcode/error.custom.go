package errorcode

import (
	errorsi "github.com/liov/hoper/go/v2/utils/errors"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type statusError ErrRep

type DefaultErrRep interface {
	ErrRep() *ErrRep
}

type GRPCStatus interface {
	GRPCStatus() *status.Status
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

func (x ErrCode) Message(msg string) error {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) Warp(err error) error {
	return &ErrRep{Code: x, Message: err.Error()}
}

func (x ErrCode) Error() string {
	return x.String()
}

func ErrHandle(err interface{}) error {
	if e, ok := err.(*ErrRep); ok {
		return e
	}
	if e, ok := err.(*ErrCode); ok {
		return e.ErrRep()
	}
	if e, ok := err.(*status.Status); ok {
		return e.Err()
	}
	if e, ok := err.(error); ok {
		return Unknown.Message(e.Error())
	}
	return Unknown.ErrRep()
}

func (x ErrCode) Origin() errorsi.ErrCode {
	return errorsi.ErrCode(x)
}

func (x ErrCode) OriErrRep() *errorsi.ErrRep {
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: x.String()}
}

func (x ErrCode) OriMessage(msg string) error {
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: msg}
}

func (x ErrCode) OriWarp(err error) error {
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: err.Error()}
}

func (x ErrCode) OriLog(err error) error {
	log.Error(err)
	return &errorsi.ErrRep{Code: errorsi.ErrCode(x), Message: x.String()}
}
