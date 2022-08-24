package client

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go.uber.org/zap"
	"time"
)

type Body struct {
	Data        []byte
	ContentType ContentType
}

func NewBody(data []byte, contentType ContentType) *Body {
	return &Body{Data: data, ContentType: contentType}
}

func (b *Body) IsJson() bool {
	return b.ContentType == ContentTypeJson
}

func (b *Body) IsProtobuf() bool {
	return b.ContentType == ContentTypeProtobuf
}

type LogCallback func(url, method, auth string, reqBody, respBody *Body, status int, process time.Duration, err error)

type Logger interface {
	SetPrefix(string)
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

var defaultLog = DefaultLogger

func DefaultLogger(url, method, auth string, reqBody, respBody *Body, status int, process time.Duration, err error) {
	reqField, respField := zap.Skip(), zap.Skip()
	if reqBody != nil {
		key := "param"
		if reqBody.IsJson() {
			reqField = zap.Reflect(key, log.BytesJson(reqBody.Data))
		} else if reqBody.IsProtobuf() {
			reqField = zap.Binary(key, reqBody.Data)
		} else {
			reqField = zap.String(key, stringsi.ToString(reqBody.Data))
		}
	}
	if respBody != nil {
		key := "result"
		if respBody.IsJson() {
			respField = zap.Reflect(key, log.BytesJson(respBody.Data))
		} else if respBody.IsProtobuf() {
			respField = zap.Binary(key, respBody.Data)
		} else {
			respField = zap.String(key, stringsi.ToString(respBody.Data))
		}
	}

	log.Default.Logger.Info("third-request", zap.String("interface", url),
		zap.String("method", method),
		reqField,
		zap.Duration("processTime", process),
		respField,
		zap.String("other", auth),
		zap.Int("status", status),
		zap.Error(err))
}
