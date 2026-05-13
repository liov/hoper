package service

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"google.golang.org/protobuf/proto"
)

const (
	RoleViewer = "viewer"
	RoleAgent  = "agent"
	maxPending = 128
)

type peer struct {
	id   string
	role string
	conn *websocket.Conn
}

type room struct {
	mu            sync.Mutex
	viewer, agent *peer
	pendToViewer  [][]byte
	pendToAgent   [][]byte
}

// RemoteBrowseHub 管理房间与信令转发；中继 TCP 由 relay.Hub 承担。
type RemoteBrowseHub struct {
	mu    sync.Mutex
	rooms map[string]*room
	// RelayTCPAddr 若非空，在双方到齐后下发给两端（host:port）。
	RelayTCPAddr string
}

func NewRemoteBrowseHub() *RemoteBrowseHub {
	return &RemoteBrowseHub{rooms: make(map[string]*room)}
}

func parseRelayAddr(addr string) (host string, port uint32) {
	host, ps, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, 19090
	}
	p, err := strconv.ParseUint(ps, 10, 32)
	if err != nil {
		return host, 19090
	}
	return host, uint32(p)
}

func (h *RemoteBrowseHub) getRoom(code string) *room {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[code]
	if !ok {
		r = &room{}
		h.rooms[code] = r
	}
	return r
}

func (h *RemoteBrowseHub) removePeerFromRoom(code, peerID string) {
	h.mu.Lock()
	r, ok := h.rooms[code]
	if !ok {
		h.mu.Unlock()
		return
	}
	h.mu.Unlock()
	r.mu.Lock()
	if r.viewer != nil && r.viewer.id == peerID {
		r.viewer = nil
	}
	if r.agent != nil && r.agent.id == peerID {
		r.agent = nil
	}
	empty := r.viewer == nil && r.agent == nil
	r.mu.Unlock()
	if empty {
		h.mu.Lock()
		delete(h.rooms, code)
		h.mu.Unlock()
	}
}

func (r *room) appendPendingViewerLocked(b []byte) {
	if len(r.pendToViewer) >= maxPending {
		r.pendToViewer = r.pendToViewer[1:]
	}
	r.pendToViewer = append(r.pendToViewer, append([]byte(nil), b...))
}

func (r *room) appendPendingAgentLocked(b []byte) {
	if len(r.pendToAgent) >= maxPending {
		r.pendToAgent = r.pendToAgent[1:]
	}
	r.pendToAgent = append(r.pendToAgent, append([]byte(nil), b...))
}

func flushPending(conn *websocket.Conn, pending *[][]byte) {
	if conn == nil {
		return
	}
	for _, b := range *pending {
		_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_ = conn.WriteMessage(websocket.BinaryMessage, b)
	}
	*pending = nil
}

func (h *RemoteBrowseHub) HandleWS(conn *websocket.Conn) {
	defer conn.Close()
	var roomCode, peerID, role string
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if roomCode != "" && peerID != "" {
				h.removePeerFromRoom(roomCode, peerID)
			}
			return
		}
		var env pb.SignalEnvelope
		if err := proto.Unmarshal(data, &env); err != nil {
			continue
		}
		switch env.Payload.(type) {
		case *pb.SignalEnvelope_Register:
			h.onRegister(conn, env.GetRegister(), &roomCode, &peerID, &role)
		case *pb.SignalEnvelope_IceParameters, *pb.SignalEnvelope_IceCandidate, *pb.SignalEnvelope_IceComplete:
			h.forwardICE(roomCode, role, &env)
		}
	}
}

func (h *RemoteBrowseHub) onRegister(conn *websocket.Conn, req *pb.RegisterReq, roomCode, peerID, role *string) {
	if req.GetRoomCode() == "" || (req.GetRole() != RoleViewer && req.GetRole() != RoleAgent) {
		_ = writeProto(conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_Error{Error: "bad register"}})
		return
	}
	*roomCode = req.GetRoomCode()
	*role = req.GetRole()
	*peerID = uuid.New().String()
	r := h.getRoom(*roomCode)
	r.mu.Lock()
	p := &peer{id: *peerID, role: *role, conn: conn}
	if *role == RoleViewer {
		if r.viewer != nil {
			r.mu.Unlock()
			_ = writeProto(conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_Error{Error: "viewer busy"}})
			return
		}
		r.viewer = p
	} else if r.agent != nil {
		r.mu.Unlock()
		_ = writeProto(conn, &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_Error{Error: "agent busy"}})
		return
	} else {
		r.agent = p
	}
	v, a := r.viewer, r.agent
	roleCopy := *role
	r.mu.Unlock()
	if roleCopy == RoleViewer {
		r.mu.Lock()
		flushPending(conn, &r.pendToViewer)
		r.mu.Unlock()
	} else {
		r.mu.Lock()
		flushPending(conn, &r.pendToAgent)
		r.mu.Unlock()
	}
	ack := &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_RegisterAck{RegisterAck: &pb.RegisterResp{
		PeerId: *peerID,
		RoomId: *roomCode,
	}}}
	_ = writeProto(conn, ack)
	if v != nil && a != nil && h.RelayTCPAddr != "" {
		sid := uuid.New().String()
		rhost, rport := parseRelayAddr(h.RelayTCPAddr)
		tok := &pb.RelayToken{SessionId: sid, RelayHost: rhost, RelayPort: rport, Psk: nil}
		envTok := &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_RelayToken{RelayToken: tok}}
		_ = writeProto(v.conn, envTok)
		_ = writeProto(a.conn, envTok)
	}
}

func (h *RemoteBrowseHub) forwardICE(roomCode, role string, env *pb.SignalEnvelope) {
	b, err := proto.Marshal(env)
	if err != nil {
		return
	}
	if roomCode == "" || role == "" {
		return
	}
	other := RoleAgent
	if role == RoleAgent {
		other = RoleViewer
	}
	r := h.getRoom(roomCode)
	r.mu.Lock()
	var target *peer
	if other == RoleViewer {
		target = r.viewer
	} else {
		target = r.agent
	}
	if target != nil {
		r.mu.Unlock()
		_ = writeBytes(target.conn, b)
		return
	}
	if other == RoleViewer {
		r.appendPendingViewerLocked(b)
	} else {
		r.appendPendingAgentLocked(b)
	}
	r.mu.Unlock()
}

func writeBytes(conn *websocket.Conn, b []byte) error {
	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteMessage(websocket.BinaryMessage, b)
}

func writeProto(conn *websocket.Conn, m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return writeBytes(conn, b)
}