package rbquic

import (
	"github.com/liov/hoper/server/go/webrtc/wire"
	"github.com/quic-go/quic-go"
)

// WriteFrame 在 QUIC stream 上写 wire 帧。
func WriteFrame(s quic.Stream, typ byte, payload []byte) error {
	return wire.WriteFrame(&s, typ, payload)
}

// ReadFrame 从 QUIC stream 读 wire 帧。
func ReadFrame(s quic.Stream) (typ byte, payload []byte, err error) {
	return wire.ReadFrame(&s)
}
