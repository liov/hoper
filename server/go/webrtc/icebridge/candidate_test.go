package icebridge_test

import (
	"testing"

	"github.com/pion/ice/v4"
	"github.com/liov/hoper/server/go/webrtc/icebridge"
)

func TestCandidateRoundTrip(t *testing.T) {
	line := "candidate:792212579 1 udp 2130706431 192.168.0.1 9 typ host"
	c1, err := ice.UnmarshalCandidate(line)
	if err != nil {
		t.Fatal(err)
	}
	p := icebridge.CandidateToProto(c1)
	c2, err := icebridge.CandidateFromProto(p)
	if err != nil {
		t.Fatal(err)
	}
	if c1.Marshal() != c2.Marshal() {
		t.Fatalf("marshal mismatch: %q vs %q", c1.Marshal(), c2.Marshal())
	}
}
