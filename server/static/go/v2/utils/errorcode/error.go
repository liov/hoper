package errorcode

import (
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
)

type statusError response.ErrorCode

func (se *statusError) Error() string {
	p := (*response.ErrorCode)(se)
	return p.Message
}

func (e ErrCode) Err() error {
	return &statusError{Code: e, Message: e.Error()}
}

func (e ErrCode) WithMessage(msg string) error {
	return &statusError{Code: e, Message: msg}
}

var SysErr = []byte(`{"code":10000,
"message":"系统错误"
}`)
