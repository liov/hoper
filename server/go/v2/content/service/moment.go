package service

import (
	"context"
	"net/http"

	"github.com/liov/hoper/go/v2/content/dao"
	model "github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MomentService struct {
	model.UnimplementedMomentServiceServer
}

func (m *MomentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "瞬间相关", "/api/moment", nil
}

func (m *MomentService) name() string {
	return "MomentService."
}

func (*MomentService) Info(context.Context, *model.GetMomentReq) (*model.Moment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthInfo not implemented")
}

func (m *MomentService) Add(ctx context.Context, req *model.AddMomentReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan(m.name() + "Add")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(Auth)
	if err != nil {
		return nil, err
	}
	log.Info(auth)
	req.UserId = auth.Id
	if err = dao.Dao.GORMDB.Save(req).Error; err != nil {
		return nil, errorcode.DBError
	}
	return nil, nil
}
func (*MomentService) Edit(context.Context, *model.AddMomentReq) (*request.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(context.Context, *model.MomentListReq) (*model.MomentListRep, error) {

	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}