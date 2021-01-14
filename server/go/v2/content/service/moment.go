package service

import (
	"context"

	model "github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MomentService struct {
	model.UnimplementedMomentServiceServer
}

func (*MomentService) Info(context.Context, *model.GetMomentReq) (*response.TinyRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (*MomentService) Add(context.Context, *model.AddMomentReq) (*response.TinyRep, error) {

	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (*MomentService) Edit(context.Context, *model.AddMomentReq) (*response.TinyRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}