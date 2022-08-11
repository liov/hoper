package output

import (
	"net/http"
	"net/url"
	"os"

	"go.uber.org/zap"
)

type httpSink struct {
	req *http.Request
}

func (*httpSink) Write(b []byte) (n int, err error) {
	return
}

func (*httpSink) Sync() error {
	return nil
}

func (*httpSink) Close() error {
	return nil
}

func RegisterSink() {
	_ = zap.RegisterSink("windows", func(u *url.URL) (sink zap.Sink, e error) {
		return os.OpenFile(u.Host+u.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	})
	_ = zap.RegisterSink("http", func(url *url.URL) (sink zap.Sink, e error) {
		s := new(httpSink)
		s.req, _ = http.NewRequest(http.MethodPost, url.Host+url.Path, nil)
		return s, nil
	})
	_ = zap.RegisterSink("https", func(url *url.URL) (sink zap.Sink, e error) {
		s := new(httpSink)
		s.req, _ = http.NewRequest(http.MethodPost, url.Host+url.Path, nil)
		return s, nil
	})
	// TODO
	_ = zap.RegisterSink("socket", func(url *url.URL) (sink zap.Sink, e error) {
		return nil, nil
	})

}
