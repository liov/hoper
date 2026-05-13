package icebridge

import (
	"context"
	"fmt"
	"sync"

	"github.com/pion/ice/v4"
	"github.com/pion/stun/v3"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
)

// SignalIO 交换 ICE 信令。
type SignalIO interface {
	SendIceParameters(ufrag, pwd string) error
	SendIceCandidate(cand *pb.IceCandidateInit) error
	SendIceComplete() error
	Recv(context.Context) (*pb.SignalEnvelope, error)
}

// Connect 完成 ICE 协商并返回数据面连接；dialer 为 true 时主动 Dial。
func Connect(ctx context.Context, urls []*stun.URI, dialer bool, sig SignalIO) (*ice.Conn, error) {
	ufrag, pwd, err := NewCredentials()
	if err != nil {
		return nil, err
	}
	ag, err := NewAgent(ctx, urls, ufrag, pwd)
	if err != nil {
		return nil, err
	}
	defer func() {
		if ag != nil {
			_ = ag.Close()
		}
	}()
	if err := bindCandidates(ag, sig); err != nil {
		return nil, err
	}
	if err := ag.GatherCandidates(); err != nil {
		return nil, err
	}
	lu, lp, err := ag.GetLocalUserCredentials()
	if err != nil {
		return nil, err
	}
	if err := sig.SendIceParameters(lu, lp); err != nil {
		return nil, err
	}
	conn, err := waitPeer(ctx, ag, dialer, sig)
	if err != nil {
		return nil, err
	}
	ag = nil
	return conn, nil
}

func bindCandidates(ag *ice.Agent, sig SignalIO) error {
	return ag.OnCandidate(func(c ice.Candidate) {
		if c == nil {
			_ = sig.SendIceComplete()
			return
		}
		_ = sig.SendIceCandidate(CandidateToProto(c))
	})
}

func waitPeer(ctx context.Context, ag *ice.Agent, dialer bool, sig SignalIO) (*ice.Conn, error) {
	var remoteUfrag, remotePwd string
	var conn *ice.Conn
	var connErr error
	var once sync.Once
	tryDial := func() {
		once.Do(func() {
			if remoteUfrag == "" {
				return
			}
			if dialer {
				conn, connErr = ag.Dial(ctx, remoteUfrag, remotePwd)
			} else {
				conn, connErr = ag.Accept(ctx, remoteUfrag, remotePwd)
			}
		})
	}
	for conn == nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if connErr != nil {
			return nil, connErr
		}
		env, err := sig.Recv(ctx)
		if err != nil {
			return nil, err
		}
		if p := env.GetIceParameters(); p != nil {
			remoteUfrag, remotePwd = p.GetUfrag(), p.GetPwd()
			_ = ag.SetRemoteCredentials(remoteUfrag, remotePwd)
			tryDial()
		}
		if c := env.GetIceCandidate(); c != nil {
			cand, err := CandidateFromProto(c)
			if err != nil {
				return nil, err
			}
			if err := ag.AddRemoteCandidate(cand); err != nil {
				return nil, err
			}
			tryDial()
		}
		if env.GetIceComplete() {
			tryDial()
		}
		if conn != nil {
			break
		}
	}
	if conn == nil {
		return nil, fmt.Errorf("icebridge: connect timeout")
	}
	ag = nil
	return conn, nil
}
