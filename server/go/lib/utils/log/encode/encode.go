package encode

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"io"
)

func DefaultReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := json.NewEncoder(w)
	// For consistency with our custom JSON encoder.
	enc.SetEscapeHTML(false)
	return enc
}
