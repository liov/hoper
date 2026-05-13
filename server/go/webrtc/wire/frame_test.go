package wire_test

import (
	"bytes"
	"testing"

	"github.com/liov/hoper/server/go/webrtc/wire"
)

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	payload := []byte("hello")
	if err := wire.WriteFrame(&buf, wire.TypeHello, payload); err != nil {
		t.Fatal(err)
	}
	typ, got, err := wire.ReadFrame(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if typ != wire.TypeHello || string(got) != "hello" {
		t.Fatalf("got typ=%d payload=%q", typ, got)
	}
}
