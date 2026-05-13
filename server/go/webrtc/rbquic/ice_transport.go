package rbquic

import (
	"context"
	"net"
	"time"

	"github.com/pion/ice/v4"
	"github.com/quic-go/quic-go"
)

type icePacketConn struct {
	c *ice.Conn
}

func (p *icePacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, err := p.c.Read(b)
	return n, p.c.RemoteAddr(), err
}

func (p *icePacketConn) WriteTo(b []byte, _ net.Addr) (int, error) {
	return p.c.Write(b)
}

func (p *icePacketConn) Close() error                       { return p.c.Close() }
func (p *icePacketConn) LocalAddr() net.Addr                { return p.c.LocalAddr() }
func (p *icePacketConn) SetDeadline(t time.Time) error      { return p.c.SetDeadline(t) }
func (p *icePacketConn) SetReadDeadline(t time.Time) error  { return p.c.SetReadDeadline(t) }
func (p *icePacketConn) SetWriteDeadline(t time.Time) error { return p.c.SetWriteDeadline(t) }

// ListenICE 在 ICE 连接上监听 QUIC。
func ListenICE(ic *ice.Conn) (*quic.Listener, error) {
	tr := &quic.Transport{Conn: &icePacketConn{c: ic}}
	return tr.Listen(DevServerTLS(), &quic.Config{})
}

// DialICE 在 ICE 连接上拨号 QUIC。
func DialICE(ctx context.Context, ic *ice.Conn) (*quic.Conn, error) {
	tr := &quic.Transport{Conn: &icePacketConn{c: ic}}
	return tr.Dial(ctx, ic.RemoteAddr(), DevClientTLS(), &quic.Config{})
}
