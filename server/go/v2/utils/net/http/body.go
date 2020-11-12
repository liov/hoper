package http2

import (
	"io"
	"net/http"
)

type NoCloseBody struct {
	s        string
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}

func (r *NoCloseBody) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

func (r *NoCloseBody) Size() int64 { return int64(len(r.s)) }

func (r *NoCloseBody) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

func (r *NoCloseBody) Close() error {
	r.i = 0
	return nil
}

//适用于轮询
func NewNoCloseBody(s string) *NoCloseBody { return &NoCloseBody{s, 0, -1} }

func NewNoCloseRequest(req *http.Request, s string) {
	v := NewNoCloseBody(s)
	req.ContentLength = int64(v.Len())
	req.Body = v
	req.GetBody = func() (io.ReadCloser, error) {
		return v, nil
	}
}
