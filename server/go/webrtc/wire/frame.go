// Package wire 定义 QUIC stream 上的轻量帧，与架构文档第 6 节一致。
package wire

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	Version1 = byte(1)
)

const (
	TypeHello = byte(iota + 1)
	TypeFileIndex
	TypeThumbReq
	TypeThumbData
	TypeFileChunk
	TypeAck
	TypeWindowUpdate
)

const HeaderLen = 6

var ErrFrameTooLarge = errors.New("wire: frame payload exceeds max")

// MaxPayload 单帧最大负载，防止 OOM（可按业务调大）。
const MaxPayload = 16 << 20

func WriteFrame(w io.Writer, typ byte, payload []byte) error {
	if len(payload) > MaxPayload {
		return ErrFrameTooLarge
	}
	var hdr [HeaderLen]byte
	hdr[0] = Version1
	hdr[1] = typ
	binary.BigEndian.PutUint32(hdr[2:], uint32(len(payload)))
	if _, err := w.Write(hdr[:]); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	_, err := w.Write(payload)
	return err
}

func ReadFrame(r io.Reader) (typ byte, payload []byte, err error) {
	var hdr [HeaderLen]byte
	if _, err = io.ReadFull(r, hdr[:]); err != nil {
		return 0, nil, err
	}
	if hdr[0] != Version1 {
		return 0, nil, errors.New("wire: bad version")
	}
	n := binary.BigEndian.Uint32(hdr[2:])
	if n > MaxPayload {
		return 0, nil, ErrFrameTooLarge
	}
	typ = hdr[1]
	if n == 0 {
		return typ, nil, nil
	}
	payload = make([]byte, n)
	if _, err = io.ReadFull(r, payload); err != nil {
		return 0, nil, err
	}
	return typ, payload, nil
}
