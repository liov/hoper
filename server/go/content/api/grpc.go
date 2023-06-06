package route

import (
	contentService "github.com/actliboy/hoper/server/go/content/service"
	"github.com/actliboy/hoper/server/go/protobuf/content"
	"google.golang.org/grpc"
)

func GrpcRegister(gs *grpc.Server) {
	content.RegisterMomentServiceServer(gs, contentService.GetMomentService())
	content.RegisterContentServiceServer(gs, contentService.GetContentService())
	content.RegisterActionServiceServer(gs, contentService.GetActionService())
}
