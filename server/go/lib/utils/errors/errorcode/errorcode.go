package errorcode

type ErrCodeInterface interface {
	ErrCode() int
	error
}

type ErrCode struct {
	Code    int
	Message string
}

func (x *ErrCode) ErrCode() int {
	return x.Code
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
