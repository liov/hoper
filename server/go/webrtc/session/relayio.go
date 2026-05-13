package session

import (
	"bytes"
	"io"
	"net"

	"github.com/liov/hoper/server/go/webrtc/relay"
	"github.com/liov/hoper/server/go/webrtc/wire"
)

func readWireRelay(c net.Conn) (typ byte, payload []byte, err error) {
	raw, err := relay.RelayRead(c)
	if err != nil {
		return 0, nil, err
	}
	return wire.ReadFrame(bytes.NewReader(raw))
}

func writeWireRelay(c net.Conn, typ byte, payload []byte) error {
	var buf bytes.Buffer
	if err := wire.WriteFrame(&buf, typ, payload); err != nil {
		return err
	}
	return relay.RelayWrite(c, buf.Bytes())
}

func readWireRelayEOF(c net.Conn) (typ byte, payload []byte, err error) {
	typ, payload, err = readWireRelay(c)
	if err == io.EOF {
		return 0, nil, err
	}
	return typ, payload, err
}
