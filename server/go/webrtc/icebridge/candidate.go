package icebridge

import (
	"fmt"
	"strings"

	"github.com/pion/ice/v4"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
)

// CandidateToProto 将 pion ICE Candidate 转为信令消息（candidate 串不含 "candidate:" 前缀也可被 UnmarshalCandidate 解析，见 pion 测试）。
func CandidateToProto(c ice.Candidate) *pb.IceCandidateInit {
	raw := c.Marshal()
	raw = strings.TrimPrefix(raw, "candidate:")
	return &pb.IceCandidateInit{
		Candidate:          raw,
		SdpMid:             "",
		SdpMlineIndex:      0,
		UsernameFragment:   "",
	}
}

// CandidateFromProto 解析为 ICE Candidate，供 Agent.AddRemoteCandidate。
func CandidateFromProto(m *pb.IceCandidateInit) (ice.Candidate, error) {
	if m == nil || m.Candidate == "" {
		return nil, fmt.Errorf("icebridge: empty candidate")
	}
	raw := m.Candidate
	if !strings.HasPrefix(raw, "candidate:") {
		raw = "candidate:" + raw
	}
	return ice.UnmarshalCandidate(raw)
}
