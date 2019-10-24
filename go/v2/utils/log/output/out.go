package output

import (
	"net/url"

	"go.uber.org/zap"
)

func RegisterSink()  {
	zap.RegisterSink("http", func(url *url.URL) (sink zap.Sink, e error) {
		return nil, nil
	})
	zap.RegisterSink("https", func(url *url.URL) (sink zap.Sink, e error) {

		return nil, nil
	})
	zap.RegisterSink("socket", func(url *url.URL) (sink zap.Sink, e error) {

		return nil, nil
	})
}


