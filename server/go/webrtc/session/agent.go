package session

import (
	"context"
	"io"
	"net"
	"strconv"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/relay"
	"github.com/liov/hoper/server/go/webrtc/rfvclient"
	"github.com/liov/hoper/server/go/webrtc/wire"
	"google.golang.org/protobuf/proto"
)

// ServeAgentRelay 在中继 agent 侧处理 wire 帧。
func ServeAgentRelay(conn net.Conn, root string) error {
	defer conn.Close()
	for {
		typ, payload, err := readWireRelay(conn)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := handleAgentRelayFrame(conn, root, typ, payload); err != nil {
			return err
		}
	}
}

func handleAgentRelayFrame(conn net.Conn, root string, typ byte, payload []byte) error {
	switch typ {
	case wire.TypeFileIndex:
		return replyFileIndexRelay(conn, root, payload)
	case wire.TypeThumbReq:
		return replyThumbRelay(conn, payload)
	default:
		return nil
	}
}

func replyFileIndexRelay(conn net.Conn, root string, payload []byte) error {
	b, err := buildFileIndexPayload(root, payload)
	if err != nil {
		return err
	}
	return writeWireRelay(conn, wire.TypeFileIndex, b)
}

func replyThumbRelay(conn net.Conn, payload []byte) error {
	b, err := buildThumbPayload(payload)
	if err != nil {
		return err
	}
	return writeWireRelay(conn, wire.TypeThumbData, b)
}

func handleAgentFrame(w io.Writer, root string, typ byte, payload []byte) error {
	switch typ {
	case wire.TypeFileIndex:
		return replyFileIndexWire(w, root, payload)
	case wire.TypeThumbReq:
		return replyThumbWire(w, payload)
	default:
		return nil
	}
}

func replyFileIndexWire(w io.Writer, root string, payload []byte) error {
	b, err := buildFileIndexPayload(root, payload)
	if err != nil {
		return err
	}
	return writeWire(w, wire.TypeFileIndex, b)
}

func replyThumbWire(w io.Writer, payload []byte) error {
	b, err := buildThumbPayload(payload)
	if err != nil {
		return err
	}
	return writeWire(w, wire.TypeThumbData, b)
}

func buildFileIndexPayload(root string, payload []byte) ([]byte, error) {
	req := &pb.ListFilesRequest{}
	if len(payload) > 0 {
		_ = proto.Unmarshal(payload, req)
	}
	path := req.GetRootPath()
	if path == "" {
		path = root
	}
	entries, err := listEntries(context.Background(), path)
	if err != nil {
		return nil, err
	}
	b, err := proto.Marshal(&pb.ListFilesResponse{Entries: entries})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func buildThumbPayload(payload []byte) ([]byte, error) {
	req := &pb.ThumbnailRequest{}
	if err := proto.Unmarshal(payload, req); err != nil {
		return nil, err
	}
	maxEdge := req.GetMaxEdge()
	if maxEdge == 0 {
		maxEdge = 256
	}
	data, hash, err := rfvclient.GetThumbnail(context.Background(), req.GetPath(), maxEdge)
	if err != nil {
		return nil, err
	}
	resp := &pb.ThumbnailResponse{Data: data, Mime: "image/webp", ThumbHash: hash}
	b, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// DialRelayAgent 连接中继并声明 agent 角色。
func DialRelayAgent(host string, port uint32, sessionID string) (net.Conn, error) {
	addr := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	if err := relay.WriteJoinHeader(c, sessionID, relay.RoleAgent); err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}

// DialRelayViewer 连接中继并声明 viewer 角色。
func DialRelayViewer(host string, port uint32, sessionID string) (net.Conn, error) {
	addr := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	if err := relay.WriteJoinHeader(c, sessionID, relay.RoleViewer); err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}
