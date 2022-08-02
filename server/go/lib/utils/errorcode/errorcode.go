package errorcode

type ErrCodeInterface interface {
	Code() int
	Error() string
}

type ErrCode struct {
	Code    int
	Message string
}

func (x *ErrCode) Error() string {
	return x.Message
}

func NewErrCode(code int, message string) *ErrCode {
	return &ErrCode{
		Code:    code,
		Message: message,
	}
}
