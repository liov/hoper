package session

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pion/ice/v4"
	"github.com/liov/hoper/server/go/webrtc/icebridge"
	"github.com/liov/hoper/server/go/webrtc/signalclient"
)

func iceTimeout() time.Duration {
	if s := strings.TrimSpace(os.Getenv("RB_ICE_TIMEOUT")); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return 12 * time.Second
}

func stunList() []string {
	raw := strings.TrimSpace(os.Getenv("RB_STUN"))
	if raw == "" {
		return []string{"stun:stun.l.google.com:19302"}
	}
	return strings.Split(raw, ",")
}

func tryICE(ctx context.Context, sess *signalclient.Session, dialer bool) (*ice.Conn, error) {
	urls, err := icebridge.ParseSTUNURIs(stunList())
	if err != nil {
		return nil, err
	}
	iceCtx, cancel := context.WithTimeout(ctx, iceTimeout())
	defer cancel()
	return icebridge.Connect(iceCtx, urls, dialer, sess)
}

// RunAgent 注册 agent，优先 ICE，失败走中继。
func RunAgent(ctx context.Context, wsURL, room, root string) error {
	cli, err := signalclient.Dial(ctx, wsURL)
	if err != nil {
		return err
	}
	defer cli.Close()
	if _, err := cli.Register(room, "agent"); err != nil {
		return err
	}
	sess := cli.StartSession()
	sess.Pump(ctx)
	iceCh := make(chan *ice.Conn, 1)
	go func() {
		c, err := tryICE(ctx, sess, false)
		if err == nil {
			iceCh <- c
		}
	}()
	select {
	case ic := <-iceCh:
		return ServeAgentICE(ic, root)
	case <-time.After(iceTimeout()):
	}
	tok, err := sess.WaitRelayToken(ctx)
	if err != nil {
		return err
	}
	conn, err := DialRelayAgent(tok.GetRelayHost(), tok.GetRelayPort(), tok.GetSessionId())
	if err != nil {
		return err
	}
	return ServeAgentRelay(conn, root)
}

// ConnectViewer 注册 viewer，优先 ICE，失败走中继 TCP。
func ConnectViewer(ctx context.Context, wsURL, room string) (io.ReadWriteCloser, error) {
	cli, err := signalclient.Dial(ctx, wsURL)
	if err != nil {
		return nil, err
	}
	if _, err := cli.Register(room, "viewer"); err != nil {
		_ = cli.Close()
		return nil, err
	}
	sess := cli.StartSession()
	sess.Pump(ctx)
	iceCh := make(chan *ice.Conn, 1)
	go func() {
		c, err := tryICE(ctx, sess, true)
		if err == nil {
			iceCh <- c
		}
	}()
	select {
	case ic := <-iceCh:
		link, err := UpgradeICEQUIC(ctx, ic, true)
		if err != nil {
			_ = cli.Close()
			return nil, err
		}
		_ = cli.Close()
		return link, nil
	case <-time.After(iceTimeout()):
	}
	tok, err := sess.WaitRelayToken(ctx)
	if err != nil {
		_ = cli.Close()
		return nil, err
	}
	conn, err := DialRelayViewer(tok.GetRelayHost(), tok.GetRelayPort(), tok.GetSessionId())
	if err != nil {
		_ = cli.Close()
		return nil, err
	}
	_ = cli.Close()
	return conn, nil
}

// ConnectViewerRelay 兼容旧调用。
func ConnectViewerRelay(ctx context.Context, wsURL, room string) (net.Conn, error) {
	link, err := ConnectViewer(ctx, wsURL, room)
	if err != nil {
		return nil, err
	}
	c, ok := link.(net.Conn)
	if !ok {
		_ = link.Close()
		return nil, fmt.Errorf("session: viewer link is not tcp relay")
	}
	return c, nil
}
