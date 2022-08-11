package log

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
)

type BytesJson json.RawMessage

func (b BytesJson) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (b BytesJson) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	return b, nil
}
func (b *BytesJson) UnmarshalJSON(raw []byte) error {
	*b = raw
	return nil
}
