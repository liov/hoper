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

func rawICEWire() bool {
	return strings.TrimSpace(os.Getenv("RB_RAW_ICE")) == "1"
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

// ViewerDirectOpt 可选的手动 IP 直连。
type ViewerDirectOpt struct {
	Host string
	Port uint32
}

// RunAgent 注册 agent，优先 IP 直连，再 ICE，失败走中继。
func RunAgent(ctx context.Context, wsURL, room, root string) error {
	ln, port, err := listenDirectTCP()
	if err != nil {
		return err
	}
	defer ln.Close()
	qln, qport, err := listenQuicRoom()
	if err != nil {
		return err
	}
	defer qln.Close()
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
	_ = sess.SendPeerEndpoints(gatherPeerEndpoints(port))
	_ = sendQuicPeerEndpoints(sess, qport)
	if c, ok := waitDirectAgent(ctx, ln, sess); ok {
		return ServeAgentWire(c, root)
	}
	if link, ok := waitQuicRoomAgent(ctx, qln); ok {
		return ServeAgentWire(link, root)
	}
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

// ConnectViewerDirect 不经信令，按 IP 直连文件端 wire。
func ConnectViewerDirect(ctx context.Context, host string, port uint32) (io.ReadWriteCloser, error) {
	c, err := dialDirectTCP(ctx, host, port)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ConnectViewer 注册 viewer，优先 IP 直连，再 ICE，失败走中继 TCP。
func ConnectViewer(ctx context.Context, wsURL, room string, opt *ViewerDirectOpt) (io.ReadWriteCloser, error) {
	cli, err := signalclient.Dial(ctx, wsURL)
	if err != nil {
		return nil, err
	}
	ln, port, err := listenDirectTCP()
	if err != nil {
		_ = cli.Close()
		return nil, err
	}
	defer ln.Close()
	if _, err := cli.Register(room, "viewer"); err != nil {
		_ = cli.Close()
		return nil, err
	}
	sess := cli.StartSession()
	sess.Pump(ctx)
	_ = sess.SendPeerEndpoints(gatherPeerEndpoints(port))
	var manualHost string
	var manualPort uint32
	if opt != nil {
		manualHost = opt.Host
		manualPort = opt.Port
	}
	if c, ok := waitDirectViewer(ctx, sess, ln, manualHost, manualPort); ok {
		_ = cli.Close()
		return c, nil
	}
	if link, ok := tryQuicRoomViewer(ctx, sess); ok {
		_ = cli.Close()
		return link, nil
	}
	iceCh := make(chan *ice.Conn, 1)
	go func() {
		c, err := tryICE(ctx, sess, true)
		if err == nil {
			iceCh <- c
		}
	}()
	select {
	case ic := <-iceCh:
		if rawICEWire() {
			_ = cli.Close()
			return ic, nil
		}
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
	link, err := ConnectViewer(ctx, wsURL, room, nil)
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
