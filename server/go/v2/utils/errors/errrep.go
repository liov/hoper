package errorsi

type ErrRep struct {
	Code    ErrCode `json:"code"`
	Message string `json:"message,omitempty"`
}
