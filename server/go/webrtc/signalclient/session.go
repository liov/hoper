package signalclient

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
)

// Session 单连接多路复用信令。
type Session struct {
	conn *websocket.Conn
	ice  chan *pb.SignalEnvelope
	peer chan *pb.PeerEndpoints
	relay chan *pb.RelayToken
	errs chan error
	once sync.Once
}

func (c *Client) StartSession() *Session {
	s := &Session{
		conn: c.conn,
		ice:  make(chan *pb.SignalEnvelope, 32),
		peer: make(chan *pb.PeerEndpoints, 8),
		relay: make(chan *pb.RelayToken, 1),
		errs: make(chan error, 1),
	}
	return s
}

func (s *Session) Pump(ctx context.Context) {
	s.once.Do(func() {
		go s.loop(ctx)
	})
}

func (s *Session) loop(ctx context.Context) {
	for {
		if err := ctx.Err(); err != nil {
			s.errs <- err
			return
		}
		env, err := readProto(s.conn)
		if err != nil {
			s.errs <- err
			return
		}
		switch env.Payload.(type) {
		case *pb.SignalEnvelope_IceParameters, *pb.SignalEnvelope_IceCandidate, *pb.SignalEnvelope_IceComplete:
			s.ice <- env
		case *pb.SignalEnvelope_PeerEndpoints:
			s.peer <- env.GetPeerEndpoints()
		case *pb.SignalEnvelope_RelayToken:
			s.relay <- env.GetRelayToken()
		default:
		}
	}
}

func (s *Session) SendIceParameters(ufrag, pwd string) error {
	return writeProto(s.conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_IceParameters{IceParameters: &pb.IceParameters{
		Ufrag: ufrag, Pwd: pwd,
	}}})
}

func (s *Session) SendIceCandidate(cand *pb.IceCandidateInit) error {
	return writeProto(s.conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_IceCandidate{IceCandidate: cand}})
}

func (s *Session) SendIceComplete() error {
	return writeProto(s.conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_IceComplete{IceComplete: true}})
}

func (s *Session) SendPeerEndpoints(eps *pb.PeerEndpoints) error {
	return writeProto(s.conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_PeerEndpoints{PeerEndpoints: eps}})
}

func (s *Session) Recv(ctx context.Context) (*pb.SignalEnvelope, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-s.errs:
		return nil, err
	case env := <-s.ice:
		return env, nil
	}
}

func (s *Session) WaitRelayToken(ctx context.Context) (*pb.RelayToken, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-s.errs:
		return nil, err
	case tok := <-s.relay:
		return tok, nil
	}
}

func (s *Session) RecvPeerEndpoints(ctx context.Context) (*pb.PeerEndpoints, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-s.errs:
		return nil, err
	case eps := <-s.peer:
		return eps, nil
	}
}
