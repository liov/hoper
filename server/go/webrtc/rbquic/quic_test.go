package rbquic_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/liov/hoper/server/go/webrtc/rbquic"
)

func TestDevListenDial(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ln, err := rbquic.ListenDev("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	addr := ln.Addr().String()
	errCh := make(chan error, 1)
	go func() {
		c, err := ln.Accept(ctx)
		if err != nil {
			errCh <- err
			return
		}
		st, err := c.AcceptStream(ctx)
		if err != nil {
			errCh <- err
			return
		}
		_, err = io.Copy(st, st)
		errCh <- err
	}()
	cli, err := rbquic.DialDev(ctx, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.CloseWithError(0, "")
	st, err := cli.OpenStreamSync(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1)
	if _, err := io.ReadFull(st, buf); err != nil {
		t.Fatal(err)
	}
}
