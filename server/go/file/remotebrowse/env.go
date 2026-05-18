package remotebrowse

import (
	"os"
	"strings"
)

// SignalUpstream 信令 WebSocket 上游（rfv-daemon）；Go 仅反代，不实现 ICE/中继逻辑。
func SignalUpstream() string {
	if u := strings.TrimSpace(os.Getenv("RB_SIGNAL_UPSTREAM")); u != "" {
		return u
	}
	return "ws://127.0.0.1:8080/rb/signal"
}

func publicSignalPath() string {
	if p := strings.TrimSpace(os.Getenv("RB_SIGNAL_WS")); p != "" {
		return p
	}
	return "/rb/signal"
}

func relayTCPHint() string {
	if a := strings.TrimSpace(os.Getenv("RB_RELAY_TCP")); a != "" {
		return a
	}
	return "see rfv-daemon (RB_RELAY_TCP)"
}
