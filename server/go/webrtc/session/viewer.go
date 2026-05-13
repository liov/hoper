package session

import (
	"io"
	"net"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/wire"
	"google.golang.org/protobuf/proto"
)

// ListFilesOn 按链路类型列举文件。
func ListFilesOn(link io.ReadWriter, root string) ([]*pb.FileEntry, error) {
	if c, ok := link.(net.Conn); ok {
		return ListFilesRelay(c, root)
	}
	return ListFilesWire(link, root)
}
func ListFilesRelay(conn net.Conn, root string) ([]*pb.FileEntry, error) {
	req := &pb.ListFilesRequest{RootPath: root}
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := writeWireRelay(conn, wire.TypeFileIndex, b); err != nil {
		return nil, err
	}
	typ, payload, err := readWireRelay(conn)
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

// FetchThumbRelay 经中继拉取缩略图。
func FetchThumbRelay(conn net.Conn, path string, maxEdge uint32) (*pb.ThumbnailResponse, error) {
	req := &pb.ThumbnailRequest{Path: path, MaxEdge: maxEdge}
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := writeWireRelay(conn, wire.TypeThumbReq, b); err != nil {
		return nil, err
	}
	typ, payload, err := readWireRelay(conn)
	if err != nil {
		return nil, err
	}
	if typ != wire.TypeThumbData {
		return nil, errUnexpectedWire(typ)
	}
	resp := &pb.ThumbnailResponse{}
	if err := proto.Unmarshal(payload, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type wireTypeErr byte

func (e wireTypeErr) Error() string {
	return "wire: unexpected frame type"
}

func errUnexpectedWire(typ byte) error {
	return wireTypeErr(typ)
}
