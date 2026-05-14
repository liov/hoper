package session

import (
	"context"
	"io"

	"github.com/pion/ice/v4"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/wire"
	"google.golang.org/protobuf/proto"
)

// ServeAgentWire 在直连数据面上处理 wire 帧。
func ServeAgentWire(rw io.ReadWriter, root string) error {
	for {
		typ, payload, err := readWire(rw)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := handleAgentFrame(rw, root, typ, payload); err != nil {
			return err
		}
	}
}

// ServeAgentICE 在 ICE 上服务 agent；RB_RAW_ICE=1 时走裸 wire。
func ServeAgentICE(conn *ice.Conn, root string) error {
	defer conn.Close()
	if rawICEWire() {
		return ServeAgentWire(conn, root)
	}
	ctx := context.Background()
	link, err := UpgradeICEQUIC(ctx, conn, false)
	if err != nil {
		return err
	}
	defer link.Close()
	return ServeAgentWire(link, root)
}

// ListFilesWire 经直连 wire 拉取列表。
func ListFilesWire(rw io.ReadWriter, root string) ([]*pb.FileEntry, error) {
	req := &pb.ListFilesRequest{RootPath: root}
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := writeWire(rw, wire.TypeFileIndex, b); err != nil {
		return nil, err
	}
	typ, payload, err := readWire(rw)
	if err != nil {
		return nil, err
	}
	if typ != wire.TypeFileIndex {
		return nil, errUnexpectedWire(typ)
	}
	resp := &pb.ListFilesResponse{}
	if err := proto.Unmarshal(payload, resp); err != nil {
		return nil, err
	}
	return resp.GetEntries(), nil
}
