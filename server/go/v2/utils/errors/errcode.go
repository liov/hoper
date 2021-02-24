package errorsi

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (x ErrCode) Message(msg string) error {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) Warp(err error) error {
	return &ErrRep{Code: x, Message: err.Error()}
}

func (x ErrCode) Log(err error) error {
	log.Error(err)
	return &ErrRep{Code: x, Message: x.String()}
}

func (x ErrCode) Error() string {
	return x.String()
}
