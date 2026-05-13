package session_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/liov/hoper/server/go/webrtc/relay"
	"github.com/liov/hoper/server/go/webrtc/session"
)

func TestRelayListRoundTrip(t *testing.T) {
	t.Setenv("RB_AGENT_USE_RFV", "0")
	h := relay.NewHub()
	ln, err := h.ServeTCP("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	addr := ln.Addr().String()
	sid := uuid.New().String()
	agentReady := make(chan struct{})
	go func() {
		host, port := splitHostPort(t, addr)
		c, err := session.DialRelayAgent(host, port, sid)
		if err != nil {
			t.Error(err)
			return
		}
		close(agentReady)
		_ = session.ServeAgentRelay(c, t.TempDir())
	}()
	time.Sleep(50 * time.Millisecond)
	host, port := splitHostPort(t, addr)
	vc, err := session.DialRelayViewer(host, port, sid)
	if err != nil {
		t.Fatal(err)
	}
	defer vc.Close()
	select {
	case <-agentReady:
	case <-time.After(3 * time.Second):
		t.Fatal("agent not ready")
	}
	entries, err := session.ListFilesRelay(vc, "")
	if err != nil {
		t.Fatal(err)
	}
	_ = entries
}
func splitHostPort(t *testing.T, addr string) (string, uint32) {
	t.Helper()
	host, ps, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	p, err := strconv.ParseUint(ps, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	return host, uint32(p)
}
