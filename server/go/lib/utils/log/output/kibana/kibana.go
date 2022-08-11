package kibana

import (
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
)

type Kibana struct {
	Es    *elasticsearch.Client
	Index string
}

func (k *Kibana) Write(b []byte) (n int, err error) {
	return
}

func (k *Kibana) Sync() error {
	return nil
}

func (k *Kibana) Close() error {
	return nil
}

// kibana://${token}?index=${index}
func RegisterSink() {
	_ = zap.RegisterSink("kibana", func(url *url.URL) (sink zap.Sink, e error) {
		k := new(Kibana)
		k.Es, e = elasticsearch.NewClient(elasticsearch.Config{})
		k.Index = url.Query().Get("index")
		return k, e
	})
}

func New(es *elasticsearch.Client, index string) zap.Sink {
	return &Kibana{es, index}
}

func NewCore() zapcore.Core { return core{} }

type core struct{}

func (core) Enabled(zapcore.Level) bool                                            { return false }
func (n core) With([]zapcore.Field) zapcore.Core                                   { return n }
func (core) Check(_ zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry { return ce }
func (core) Write(zapcore.Entry, []zapcore.Field) error                            { return nil }
func (core) Sync() error                                                           { return nil }
