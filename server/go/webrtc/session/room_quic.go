package session

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/liov/hoper/server/go/webrtc/rbquic"
	"github.com/liov/hoper/server/go/webrtc/signalclient"
	"github.com/quic-go/quic-go"
)

func quicRoomPort() int {
	if s := strings.TrimSpace(os.Getenv("RB_QUIC_PORT")); s != "" {
		if p, err := strconv.Atoi(s); err == nil && p > 0 {
			return p
		}
	}
	return 19092
}

func listenQuicRoom() (*quic.EarlyListener, uint32, error) {
	ln, err := rbquic.ListenDev(fmt.Sprintf(":%d", quicRoomPort()))
	if err != nil {
		return nil, 0, err
	}
	_, ps, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		_ = ln.Close()
		return nil, 0, err
	}
	p, err := strconv.ParseUint(ps, 10, 32)
	if err != nil {
		_ = ln.Close()
		return nil, 0, err
	}
	return ln, uint32(p), nil
}

func waitQuicRoomAgent(ctx context.Context, ln *quic.EarlyListener) (io.ReadWriteCloser, bool) {
	ctx, cancel := context.WithTimeout(ctx, directTimeout())
	defer cancel()
	qc, err := ln.Accept(ctx)
	if err != nil {
		return nil, false
	}
	st, err := qc.AcceptStream(ctx)
	if err != nil {
		_ = qc.CloseWithError(0, "stream")
		return nil, false
	}
	return &quicStream{qc: qc, st: st}, true
}

func tryQuicRoomViewer(ctx context.Context, sess *signalclient.Session) (io.ReadWriteCloser, bool) {
	eps, err := waitPeerEndpoints(ctx, sess)
	if err != nil {
		return nil, false
	}
	port := uint32(quicRoomPort())
	for _, ep := range eps.GetItems() {
		if ep.GetHost() == "" {
			continue
		}
		link, err := dialQuicRoom(ctx, ep.GetHost(), port)
		if err == nil {
			return link, true
		}
	}
	return nil, false
}

func dialQuicRoom(ctx context.Context, host string, port uint32) (io.ReadWriteCloser, error) {
	qc, err := rbquic.DialDev(ctx, net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10)))
	if err != nil {
		return nil, err
	}
	st, err := qc.OpenStreamSync(ctx)
	if err != nil {
		_ = qc.CloseWithError(0, "stream")
		return nil, err
	}
	return &quicStream{qc: qc, st: st}, nil
}

func sendQuicPeerEndpoints(sess *signalclient.Session, port uint32) error {
	return sess.SendPeerEndpoints(gatherPeerEndpoints(port))
}
