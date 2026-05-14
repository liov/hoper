package session

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/signalclient"
)

func directTimeout() time.Duration {
	if s := strings.TrimSpace(os.Getenv("RB_DIRECT_TIMEOUT")); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return 5 * time.Second
}

func directTCPPort() int {
	if s := strings.TrimSpace(os.Getenv("RB_DIRECT_TCP")); s != "" {
		if p, err := strconv.Atoi(s); err == nil && p > 0 {
			return p
		}
	}
	return 19091
}

func listenDirectTCP() (net.Listener, uint32, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", directTCPPort()))
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

func gatherPeerEndpoints(port uint32) *pb.PeerEndpoints {
	var items []*pb.PeerEndpoint
	ifaces, err := net.Interfaces()
	if err != nil {
		return &pb.PeerEndpoints{}
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ip := addrIP(a)
			if ip == nil || ip.IsLoopback() || ip.IsUnspecified() {
				continue
			}
			items = append(items, &pb.PeerEndpoint{Host: ip.String(), Port: port})
		}
	}
	return &pb.PeerEndpoints{Items: items}
}

func addrIP(a net.Addr) net.IP {
	switch v := a.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		return nil
	}
}

func dialDirectTCP(ctx context.Context, host string, port uint32) (net.Conn, error) {
	d := net.Dialer{Timeout: directTimeout()}
	return d.DialContext(ctx, "tcp", net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10)))
}

func tryDialPeerEndpoints(ctx context.Context, eps *pb.PeerEndpoints) (net.Conn, error) {
	for _, ep := range eps.GetItems() {
		if ep.GetHost() == "" || ep.GetPort() == 0 {
			continue
		}
		c, err := dialDirectTCP(ctx, ep.GetHost(), ep.GetPort())
		if err == nil {
			return c, nil
		}
	}
	return nil, fmt.Errorf("direct: no reachable endpoint")
}

func waitPeerEndpoints(ctx context.Context, sess *signalclient.Session) (*pb.PeerEndpoints, error) {
	for {
		eps, err := sess.RecvPeerEndpoints(ctx)
		if err != nil {
			return nil, err
		}
		if eps != nil && len(eps.GetItems()) > 0 {
			return eps, nil
		}
	}
}
