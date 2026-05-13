package service

import (
	"os"
	"sync"

	"github.com/liov/hoper/server/go/webrtc/relay"
)

var (
	signalOnce sync.Once
	signalHub  *RemoteBrowseHub
)

// EnsureSignalHub 启动信令 Hub 与中继监听（进程内单例）。
func EnsureSignalHub() *RemoteBrowseHub {
	signalOnce.Do(func() {
		signalHub = NewRemoteBrowseHub()
		addr := os.Getenv("RB_RELAY_TCP")
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
