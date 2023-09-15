package route

import (
	contentService "github.com/liovx/hoper/server/go/content/service"
	"github.com/liovx/hoper/server/go/protobuf/content"
	"google.golang.org/grpc"
)

func GrpcRegister(gs *grpc.Server) {
	content.RegisterMomentServiceServer(gs, contentService.GetMomentService())
	content.RegisterContentServiceServer(gs, contentService.GetContentService())
	content.RegisterActionServiceServer(gs, contentService.GetActionService())
}
