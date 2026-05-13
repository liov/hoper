// Package icebridge 封装 pion/ice（不含 WebRTC），供 Go 侧与 quic-go 共用 UDP 前的 ICE 建连。
// 控制端调用 Agent.StartDial / Dial；被控端调用 StartAccept / Accept（见 pion/ice transport）。
package icebridge

import (
	"context"
	"fmt"

	"github.com/pion/ice/v4"
	"github.com/pion/stun/v3"
)

// ParseSTUNURIs 解析 STUN/TURN URL 列表，如 stun:stun.l.google.com:19302。
func ParseSTUNURIs(urls []string) ([]*stun.URI, error) {
	out := make([]*stun.URI, 0, len(urls))
	for _, u := range urls {
		uri, err := stun.ParseURI(u)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", u, err)
		}
		out = append(out, uri)
	}
	return out, nil
}

// NewAgent 创建 ICE Agent。后续：SetRemoteCredentials → GatherCandidates → OnCandidate 交换 → AddRemoteCandidate → Dial/Accept。
func NewAgent(_ context.Context, urls []*stun.URI, localUfrag, localPwd string) (*ice.Agent, error) {
	if localUfrag == "" || localPwd == "" {
		return nil, fmt.Errorf("icebridge: empty ufrag or pwd")
	}
	return ice.NewAgent(&ice.AgentConfig{
		Urls:         urls,
		LocalUfrag:   localUfrag,
		LocalPwd:     localPwd,
		NetworkTypes: []ice.NetworkType{ice.NetworkTypeUDP6, ice.NetworkTypeUDP4},
	})
}

// DefaultSTUN 与仓库内历史示例一致的公共 STUN，生产应改为自建。
func DefaultSTUN() []*stun.URI {
	u, _ := stun.ParseURI("stun:stun.l.google.com:19302")
	return []*stun.URI{u}
}
