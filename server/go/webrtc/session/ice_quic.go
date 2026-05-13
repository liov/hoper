package session

import (
	"context"
	"io"

	"github.com/pion/ice/v4"
	"github.com/quic-go/quic-go"
	"github.com/liov/hoper/server/go/webrtc/rbquic"
)

type quicStream struct {
	qc *quic.Conn
	st *quic.Stream
}

func (q *quicStream) Read(p []byte) (int, error)  { return q.st.Read(p) }
func (q *quicStream) Write(p []byte) (int, error) { return q.st.Write(p) }
func (q *quicStream) Close() error {
	_ = q.st.Close()
	return q.qc.CloseWithError(0, "")
}

// UpgradeICEQUIC 在 ICE 上建立 QUIC 单流数据面。
func UpgradeICEQUIC(ctx context.Context, ic *ice.Conn, dialer bool) (io.ReadWriteCloser, error) {
	if dialer {
		qc, err := rbquic.DialICE(ctx, ic)
		if err != nil {
			return ic, nil
		}
		st, err := qc.OpenStreamSync(ctx)
		if err != nil {
			_ = qc.CloseWithError(0, "")
			return ic, nil
		}
		return &quicStream{qc: qc, st: st}, nil
	}
	ln, err := rbquic.ListenICE(ic)
	if err != nil {
		return ic, nil
	}
	qc, err := ln.Accept(ctx)
	if err != nil {
		return ic, nil
	}
	st, err := qc.AcceptStream(ctx)
	if err != nil {
		_ = qc.CloseWithError(0, "")
		return ic, nil
	}
	return &quicStream{qc: qc, st: st}, nil
}
