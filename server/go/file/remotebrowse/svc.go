package remotebrowse

import (
	"context"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/file/rfvclient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service 远程相册 HTTP/gRPC 门面：请求转发至 rfv gRPC，不实现 P2P/ICE。
type Service struct {
	pb.UnimplementedRemoteBrowseServiceServer
}

var browseSvc *Service

func GetService() *Service {
	if browseSvc == nil {
		browseSvc = &Service{}
	}
	return browseSvc
}

func (s *Service) GetHealth(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	_ = ctx
	return &pb.HealthResponse{
		SignalWs:   publicSignalPath(),
		RelayTcp:   relayTCPHint(),
		RfvGrpc:    rfvclient.GRPCAddr(),
		ThumbCache: rfvclient.ThumbCacheDir(),
	}, nil
}

func (s *Service) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	if req.GetRootPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing path")
	}
	entries, err := rfvclient.ListRemoteFiles(ctx, req.GetRootPath())
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "rfv: %v", err)
	}
	return &pb.ListFilesResponse{Entries: entries}, nil
}

func (s *Service) GetThumbnail(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
	if req.GetHash() != "" {
		if b, ok := rfvclient.GetThumbnailCached(req.GetHash()); ok {
			return &pb.ThumbnailResponse{Data: b, Mime: "image/webp", ThumbHash: req.GetHash()}, nil
		}
	}
	if req.GetPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing path or hash")
	}
	maxEdge := req.GetMaxEdge()
	if maxEdge == 0 {
		maxEdge = 256
	}
	data, hash, err := rfvclient.GetThumbnailStream(ctx, req.GetPath(), maxEdge)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "rfv: %v", err)
	}
	return &pb.ThumbnailResponse{Data: data, Mime: "image/webp", ThumbHash: hash}, nil
}
