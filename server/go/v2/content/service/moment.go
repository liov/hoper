package service

import (
	"context"
	"net/http"

	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/protobuf/content"
	model "github.com/liov/hoper/go/v2/protobuf/user"
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

func (*MomentService) Info(context.Context, *content.GetMomentReq) (*content.Moment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthInfo not implemented")
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Empty, error) {
	ctxi, span := model.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	user, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	err = dao.Dao.Limit(ctxi, &conf.Conf.Customize.Limit)
	if err != nil {
		return nil, err
	}
	log.Info(user)
	req.UserId = user.Id
	db := dao.Dao.GORMDB
	/*	var count int64
		db.Table(`mood`).Where(`name = ?`, req.MoodName).Count(&count)
		if count == 0 {
			return nil, errorcode.ParamInvalid.Message("心情不存在")
		}*/
	var tags []*content.Tag
	db.Table("tag").Select("id,name").
		Where("name IN (?)", req.Tags).Find(&tags)

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
