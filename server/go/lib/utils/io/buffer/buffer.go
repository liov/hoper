package buffer

import "go.uber.org/zap/buffer"

var (
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = buffer.NewPool().Get
)
