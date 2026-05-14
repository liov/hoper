package remotebrowse

import (
	"os"
	"strings"
	"sync"

	"github.com/liov/hoper/server/go/file/relay"
)

var (
	signalOnce sync.Once
	signalHub  *Hub
)

// EnsureHub 返回信令 Hub；未配置 RB_SIGNAL_UPSTREAM 时在本进程内启动中继。
func EnsureHub() *Hub {
	signalOnce.Do(func() {
		signalHub = NewHub()
		if strings.TrimSpace(os.Getenv("RB_SIGNAL_UPSTREAM")) != "" {
			signalHub.RelayTCPAddr = strings.TrimSpace(os.Getenv("RB_RELAY_TCP"))
			return
		}
		addr := strings.TrimSpace(os.Getenv("RB_RELAY_TCP"))
		if addr == "" {
			addr = "127.0.0.1:0"
		}
		if ln, err := DefaultRelayHub.ServeTCP(addr); err == nil {
			signalHub.RelayTCPAddr = ln.Addr().String()
		}
	})
	return signalHub
}

// DefaultRelayHub 全局中继（与信令同进程时可复用）。
var DefaultRelayHub = relay.NewHub()
