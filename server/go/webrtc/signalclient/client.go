package signalclient

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"google.golang.org/protobuf/proto"
)

// Client 远程浏览信令 WebSocket 客户端。
type Client struct {
	conn *websocket.Conn
}

func Dial(ctx context.Context, url string) (*Client, error) {
	d := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	c, resp, err := d.DialContext(ctx, url, http.Header{})
	if err != nil {
		return nil, err
	}
	if resp != nil {
		_ = resp.Body.Close()
	}
	return &Client{conn: c}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Register(room, role string) (*pb.RegisterResp, error) {
	env := &pb.SignalEnvelope{Payload: &pb.SignalEnvelope_Register{Register: &pb.RegisterReq{
		RoomCode: room, Role: role,
	}}}
	if err := writeProto(c.conn, env); err != nil {
		return nil, err
	}
	for {
		in, err := readProto(c.conn)
		if err != nil {
			return nil, err
		}
		if ack := in.GetRegisterAck(); ack != nil {
			return ack, nil
		}
		if msg := in.GetError(); msg != "" {
			return nil, errSignal(msg)
		}
	}
}

func (c *Client) WaitRelayToken(ctx context.Context) (*pb.RelayToken, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		in, err := readProto(c.conn)
		if err != nil {
			if ne, ok := err.(netError); ok && ne.Timeout() {
				continue
			}
			return nil, err
		}
		_ = c.conn.SetReadDeadline(time.Time{})
		if tok := in.GetRelayToken(); tok != nil {
			return tok, nil
		}
		if msg := in.GetError(); msg != "" {
			return nil, errSignal(msg)
		}
	}
}

type netError interface {
	Timeout() bool
}

type signalErr string

func (e signalErr) Error() string { return string(e) }

func errSignal(msg string) error { return signalErr(msg) }

func writeProto(conn *websocket.Conn, env *pb.SignalEnvelope) error {
	b, err := proto.Marshal(env)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.BinaryMessage, b)
}

func readProto(conn *websocket.Conn) (*pb.SignalEnvelope, error) {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var env pb.SignalEnvelope
	if err := proto.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return &env, nil
}

// Drain 忽略后续 ICE 等信令，直至连接关闭。
func (c *Client) Drain() {
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
}
