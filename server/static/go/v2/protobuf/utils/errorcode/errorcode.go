package errorcode

func (x ErrCode) Error() string {
	return x.String()
}
