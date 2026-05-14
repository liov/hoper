package session

import (
	"context"
	"net"

	"github.com/liov/hoper/server/go/webrtc/signalclient"
)

func waitDirectAgent(ctx context.Context, ln net.Listener, sess *signalclient.Session) (net.Conn, bool) {
	dctx, cancel := context.WithTimeout(ctx, directTimeout())
	defer cancel()
	ch := make(chan net.Conn, 1)
	go acceptDirectOnce(dctx, ln, ch)
	go dialPeerOnce(dctx, sess, ch)
	select {
	case c := <-ch:
		return c, true
	case <-dctx.Done():
		return nil, false
	}
}

func waitDirectViewer(ctx context.Context, sess *signalclient.Session, ln net.Listener, manualHost string, manualPort uint32) (net.Conn, bool) {
	dctx, cancel := context.WithTimeout(ctx, directTimeout())
	defer cancel()
	ch := make(chan net.Conn, 1)
	if manualHost != "" && manualPort != 0 {
		go dialManualOnce(dctx, manualHost, manualPort, ch)
	}
	if ln != nil {
		go acceptDirectOnce(dctx, ln, ch)
	}
	go dialPeerOnce(dctx, sess, ch)
	select {
	case c := <-ch:
		return c, true
	case <-dctx.Done():
		return nil, false
	}
}

func acceptDirectOnce(ctx context.Context, ln net.Listener, ch chan<- net.Conn) {
	if ln == nil {
		return
	}
	type res struct {
		c   net.Conn
		err error
	}
	done := make(chan res, 1)
	go func() {
		c, err := ln.Accept()
		done <- res{c: c, err: err}
	}()
	select {
	case <-ctx.Done():
		return
	case r := <-done:
		if r.err == nil {
			ch <- r.c
		}
	}
}

func dialManualOnce(ctx context.Context, host string, port uint32, ch chan<- net.Conn) {
	c, err := dialDirectTCP(ctx, host, port)
	if err == nil {
		ch <- c
	}
}

func dialPeerOnce(ctx context.Context, sess *signalclient.Session, ch chan<- net.Conn) {
	eps, err := waitPeerEndpoints(ctx, sess)
	if err != nil || eps == nil {
		return
	}
	c, err := tryDialPeerEndpoints(ctx, eps)
	if err == nil {
		ch <- c
	}
}
