package errorcode

import "github.com/liov/hoper/go/v2/utils/log"

type statusError ErrRep

func (se *statusError) Error() string {
	p := (*ErrRep)(se)
	return p.Message
}

func (x ErrCode) Err() error {
	return &statusError{Code: x, Message: x.Error()}
}

func (x ErrCode) WithMessage(msg string) error {
	return &statusError{Code: x, Message: msg}
}

func (x ErrCode) Log(err error) error {
	log.Default.Error(err)
	return &statusError{Code: x, Message: x.Error()}
}

var SysErr = []byte(`{"code":10000,
"message":"系统错误"
}`)
