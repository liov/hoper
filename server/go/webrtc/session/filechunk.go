package session

import (
	"context"
	"io"
	"os"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/wire"
	"google.golang.org/protobuf/proto"
)

const fileChunkSize = 256 << 10

// FetchFileChunkWire 读取原文件分片。
func FetchFileChunkWire(rw io.ReadWriter, path string, off, length int64) ([]byte, error) {
	req := &pb.OpenFileRequest{Id: path, Variant: pb.FileVariant_FILE_VARIANT_ORIGINAL, ByteOffset: off, ByteLength: length}
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := writeWire(rw, wire.TypeFileChunk, b); err != nil {
		return nil, err
	}
	typ, payload, err := readWire(rw)
	if err != nil {
		return nil, err
	}
	if typ != wire.TypeFileChunk {
		return nil, errUnexpectedWire(typ)
	}
	return payload, nil
}

func replyFileChunkWire(w io.Writer, root string, payload []byte) error {
	b, err := buildFileChunkPayload(root, payload)
	if err != nil {
		return err
	}
	return writeWire(w, wire.TypeFileChunk, b)
}

func buildFileChunkPayload(root string, payload []byte) ([]byte, error) {
	req := &pb.OpenFileRequest{}
	if err := proto.Unmarshal(payload, req); err != nil {
		return nil, err
	}
	path, err := ResolveUnderRoot(root, req.GetId())
	if err != nil {
		return nil, err
	}
	return readFileChunk(path, req.GetByteOffset(), req.GetByteLength())
}

func readFileChunk(path string, off, length int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if off > 0 {
		if _, err := f.Seek(off, io.SeekStart); err != nil {
			return nil, err
		}
	}
	n := length
	if n <= 0 || n > fileChunkSize {
		n = fileChunkSize
	}
	buf := make([]byte, n)
	got, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf[:got], nil
}

func listEntriesSafe(ctx context.Context, root, reqPath string) ([]*pb.FileEntry, error) {
	dir, err := ResolveUnderRoot(root, reqPath)
	if err != nil {
		return nil, err
	}
	return listEntries(ctx, dir)
}
