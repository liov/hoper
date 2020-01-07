package errorcode

import (
	"github.com/liov/hoper/go/v2/protobuf/utils"
)

type statusError utils.ErrorCode

func (se *statusError) Error() string {
	p := (*utils.ErrorCode)(se)
	return p.Message
}

func ErrorWithMessage(c ErrCode, msg string) error {
	return &statusError{Code: uint32(c), Message: msg}
}

func Error(c ErrCode) error {
	return &statusError{Code: uint32(c), Message: c.Error()}
}

func (e ErrCode) Err() error {
	return &statusError{Code: uint32(e), Message: e.Error()}
}

func (e ErrCode) WithMessage(msg string) error {
	return &statusError{Code: uint32(e), Message: msg}
}

func (e ErrCode) WithError(err error) error {
	return &statusError{Code: uint32(e), Message: err.Error()}
}
