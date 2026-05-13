package session

import (
	"io"

	"github.com/liov/hoper/server/go/webrtc/wire"
)

func readWire(r io.Reader) (typ byte, payload []byte, err error) {
	return wire.ReadFrame(r)
}

func writeWire(w io.Writer, typ byte, payload []byte) error {
	return wire.WriteFrame(w, typ, payload)
}
