package service

import (
	"context"
	"net/http"

	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MomentService struct {
	content.UnimplementedMomentServiceServer
}

func (m *MomentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "瞬间相关", "/api/moment", nil
}

func (m *MomentService) name() string {
	return "MomentService."
}

func (*MomentService) Info(context.Context, *content.GetMomentReq) (*content.Moment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthInfo not implemented")
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan(m.name() + "Add")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(Auth)
	if err != nil {
		return nil, err
	}
	err = Limit(ctxi, model.MomentMinuteLimitKey, model.MomentMinuteLimit, model.MomentDayLimitKey, model.MomentDayLimit)
	if err != nil {
		return nil, err
	}
	log.Info(auth)
	req.UserId = auth.Id
	db := dao.Dao.GORMDB
	var count int64
	db.Table(`mood`).Where(`name = ?`,req.MoodName).Count(&count)
	if count == 0 {
		return nil, errorcode.ParamInvalid.Message("心情不存在")
	}
	if err = db.Save(req).Error; err != nil {
		return nil, errorcode.DBError
	}
	return nil, nil
}
func (*MomentService) Edit(context.Context, *content.AddMomentReq) (*request.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(context.Context, *content.MomentListReq) (*content.MomentListRep, error) {

	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
