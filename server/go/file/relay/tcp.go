// Package relay 实现自研 TCP 中继：双方以 session 关联，仅转发密文 payload。
package relay

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	magic    = "RBRL"
	frameVer = byte(1)
)

const (
	RoleViewer = byte(0)
	RoleAgent  = byte(1)
)

var (
	ErrBadMagic = errors.New("relay: bad magic")
	ErrBadFrame = errors.New("relay: bad frame")
	ErrDupRole  = errors.New("relay: duplicate role for session")
)

type sessionPair struct {
	mu       sync.Mutex
	cond     *sync.Cond
	conns    [2]net.Conn
	bridging bool
}

func newSessionPair() *sessionPair {
	s := &sessionPair{}
	s.cond = sync.NewCond(&s.mu)
	return s
}

type Hub struct {
	mu       sync.Mutex
	sessions map[string]*sessionPair
}

func NewHub() *Hub {
	return &Hub{sessions: make(map[string]*sessionPair)}
}

func (h *Hub) session(sid string) *sessionPair {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s, ok := h.sessions[sid]; ok {
		return s
	}
	s := newSessionPair()
	h.sessions[sid] = s
	return s
}

func (h *Hub) deleteSession(sid string) {
	h.mu.Lock()
	delete(h.sessions, sid)
	h.mu.Unlock()
}

// ServeTCP 在 addr 上监听；返回 Listener 便于测试关闭。
func (h *Hub) ServeTCP(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go h.handleConn(c)
		}
	}()
	return ln, nil
}

func readJoin(r io.Reader) (sessionID string, role byte, err error) {
	var magicBuf [4]byte
	if _, err = io.ReadFull(r, magicBuf[:]); err != nil {
		return "", 0, err
	}
	if string(magicBuf[:]) != magic {
		return "", 0, ErrBadMagic
	}
	var hdr [1 + 16 + 1]byte
	if _, err = io.ReadFull(r, hdr[:]); err != nil {
		return "", 0, err
	}
	if hdr[0] != frameVer {
		return "", 0, ErrBadFrame
	}
	sessionID = uuid.UUID(hdr[1:17]).String()
	role = hdr[17]
	if role != RoleViewer && role != RoleAgent {
		return "", 0, ErrBadFrame
	}
	return sessionID, role, nil
}

// WriteJoinHeader 客户端首包：魔数 + 版本 + session uuid + role。
func WriteJoinHeader(w io.Writer, sessionID string, role byte) error {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}
	var buf [4 + 1 + 16 + 1]byte
	copy(buf[0:4], magic)
	buf[4] = frameVer
	copy(buf[5:21], id[:])
	buf[21] = role
	_, err = w.Write(buf[:])
	return err
}

func readDataFrame(r io.Reader) ([]byte, error) {
	var sz [4]byte
	if _, err := io.ReadFull(r, sz[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(sz[:])
	if n > 1<<20 {
		return nil, ErrBadFrame
	}
	if n == 0 {
		return nil, nil
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func writeDataFrame(w io.Writer, payload []byte) error {
	var sz [4]byte
	binary.BigEndian.PutUint32(sz[:], uint32(len(payload)))
	if _, err := w.Write(sz[:]); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	_, err := w.Write(payload)
	return err
}

func (h *Hub) handleConn(c net.Conn) {
	_ = c.SetDeadline(time.Now().Add(45 * time.Second))
	sessionID, role, err := readJoin(c)
	if err != nil {
		_ = c.Close()
		return
	}
	_ = c.SetDeadline(time.Time{})
	idx := int(role)
	sp := h.session(sessionID)
	sp.mu.Lock()
	if sp.conns[idx] != nil {
		sp.mu.Unlock()
		_ = c.Close()
		return
	}
	sp.conns[idx] = c
	sp.cond.Broadcast()
	other := 1 - idx
	deadline := time.Now().Add(120 * time.Second)
	for sp.conns[other] == nil {
		if time.Now().After(deadline) {
			sp.conns[idx] = nil
			sp.cond.Broadcast()
			sp.mu.Unlock()
			_ = c.Close()
			h.maybeDeleteEmpty(sessionID, sp)
			return
		}
		sp.cond.Wait()
	}
	if sp.bridging {
		sp.mu.Unlock()
		return
	}
	sp.bridging = true
	v, a := sp.conns[0], sp.conns[1]
	sp.mu.Unlock()
	go h.bridge(sessionID, v, a, sp)
}

func (h *Hub) maybeDeleteEmpty(sid string, sp *sessionPair) {
	sp.mu.Lock()
	empty := sp.conns[0] == nil && sp.conns[1] == nil
	sp.mu.Unlock()
	if empty {
		h.deleteSession(sid)
	}
}

func (h *Hub) bridge(sid string, viewer, agent net.Conn, sp *sessionPair) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		pipe(viewer, agent)
	}()
	go func() {
		defer wg.Done()
		pipe(agent, viewer)
	}()
	wg.Wait()
	_ = viewer.Close()
	_ = agent.Close()
	sp.mu.Lock()
	sp.conns[0], sp.conns[1] = nil, nil
	sp.bridging = false
	sp.mu.Unlock()
	h.deleteSession(sid)
}

func pipe(dst, src net.Conn) {
	for {
		payload, err := readDataFrame(src)
		if err != nil {
			return
		}
		if err := writeDataFrame(dst, payload); err != nil {
			return
		}
	}
}

// RelayWrite 写入一帧密文到已建立 join 的连接。
func RelayWrite(c net.Conn, payload []byte) error {
	return writeDataFrame(c, payload)
}

// RelayRead 从连接读一帧。
func RelayRead(c net.Conn) ([]byte, error) {
	return readDataFrame(c)
}
