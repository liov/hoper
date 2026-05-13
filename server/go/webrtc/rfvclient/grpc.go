// Package rfvclient 通过 gRPC 调用 rfv RemoteBrowseMedia。
package rfvclient

import (
	"context"
	"os"
	"strings"
	"sync"

	grpcx "github.com/hopeio/gox/net/http/grpc"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"google.golang.org/grpc"
)

var (
	grpcOnce sync.Once
	grpcConn *grpc.ClientConn
	grpcCli  pb.RemoteBrowseServiceClient
)

// GRPCAddr 环境变量 RFV_GRPC，默认 127.0.0.1:50051。
func GRPCAddr() string {
	s := strings.TrimSpace(os.Getenv("RFV_GRPC"))
	if s == "" {
		return "127.0.0.1:50051"
	}
	return s
}

func client() (pb.RemoteBrowseServiceClient, error) {
	var err error
	grpcOnce.Do(func() {
		grpcConn, err = grpcx.NewClient(GRPCAddr())
		if err != nil {
			return
		}
		grpcCli = pb.NewRemoteBrowseServiceClient(grpcConn)
	})
	return grpcCli, err
}

// ListRemoteFiles 列举远端目录文件元数据。
func ListRemoteFiles(ctx context.Context, rootPath string) ([]*pb.FileEntry, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	resp, err := c.ListFiles(ctx, &pb.ListFilesRequest{RootPath: rootPath})
	if err != nil {
		return nil, err
	}
	return resp.GetEntries(), nil
}

// GetThumbnail 拉取缩略图；优先读 Go 侧本地缓存（thumb_hash 命中则不再 RPC）。
func GetThumbnail(ctx context.Context, path string, maxEdge uint32) ([]byte, string, error) {
	if maxEdge == 0 {
		maxEdge = 256
	}
	c, err := client()
	if err != nil {
		return nil, "", err
	}
	resp, err := c.GetThumbnail(ctx, &pb.ThumbnailRequest{Path: path, MaxEdge: maxEdge})
	if err != nil {
		return nil, "", err
	}
	hash := resp.GetThumbHash()
	if err := StoreThumbCache(hash, resp.GetData()); err != nil {
		return nil, "", err
	}
	return resp.GetData(), hash, nil
}

// GetThumbnailCached 若本地已有 thumb_hash 对应文件则直接返回。
func GetThumbnailCached(hash string) ([]byte, bool) {
	return LoadThumbCache(hash)
}
