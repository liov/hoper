package errorcode

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type statusError ErrRep

type GRPCErr interface {
	GRPCErr() *ErrRep
}

func (m *ErrRep) Error() string {
	return m.Message
}

func (x *ErrRep) GRPCStatus() *status.Status {
	return status.New(codes.Code(x.Code), x.Message)
}

func (x ErrCode) GRPCErr() *ErrRep {
	return &ErrRep{Code: x, Message: x.String()}
}

//example 实现
func (x ErrCode) GRPCStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x ErrCode) WithMessage(msg string) error {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) Log(err error) error {
	log.Default.Error(err)
	return &ErrRep{Code: x, Message: x.String()}
}

func (x ErrCode) Error() string {
	return x.String()
}

var SysErr = []byte(`{"code":10000,
"message":"系统错误"
}`)
