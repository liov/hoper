package rfvclient

import (
	"context"
	"errors"
	"io"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
)

// GetThumbnailStream 经双向流拉取缩略图，失败回退一元 RPC。
func GetThumbnailStream(ctx context.Context, path string, maxEdge uint32) ([]byte, string, error) {
	if maxEdge == 0 {
		maxEdge = 256
	}
	c, err := client()
	if err != nil {
		return nil, "", err
	}
	stream, err := c.ThumbnailPipe(ctx)
	if err != nil {
		return GetThumbnail(ctx, path, maxEdge)
	}
	defer stream.CloseSend()
	if err := stream.Send(&pb.ThumbnailRequest{Path: path, MaxEdge: maxEdge}); err != nil {
		return GetThumbnail(ctx, path, maxEdge)
	}
	var out []byte
	var hash, mime string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return GetThumbnail(ctx, path, maxEdge)
		}
		if msg := chunk.GetError(); msg != "" {
			return nil, "", errors.New(msg)
		}
		out = append(out, chunk.GetData()...)
		hash = chunk.GetThumbHash()
		mime = chunk.GetMime()
		if chunk.GetDone() {
			break
		}
	}
	if len(out) == 0 {
		return GetThumbnail(ctx, path, maxEdge)
	}
	if hash != "" {
		_ = StoreThumbCache(hash, out)
	}
	if mime == "" {
		mime = "image/webp"
	}
	return out, hash, nil
}
