package service

import (
	"context"

	"github.com/hopeio/gox/log"
	"go.uber.org/zap"

	"github.com/hopeio/protobuf/request"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommonService struct {
	common.UnimplementedCommonServiceServer
}

func (*CommonService) TagInfo(context.Context, *request.Id) (*common.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (*CommonService) AddTag(ctx context.Context, req *common.AddTagReq) (*emptypb.Empty, error) {

	user, err := auth(ctx, false)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB
	req.UserId = user.Id
	err = db.Create(req).Error
	if err != nil {
		log.Errorw("db.Create", zap.Error(err))
		return nil, err
	}
	return nil, nil
}
func (*CommonService) EditTag(ctx context.Context, req *common.EditTagReq) (*emptypb.Empty, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}

	db := global.Dao.GORMDB.DB
	err = db.Updates(&common.Tag{
		Desc:  req.Desc,
		Image: req.Image,
	}).Where(`id = ? AND user_id = ? AND status = 0`, req.Id, auth.Id).Error
	if err != nil {
		log.Errorw("db.Updates", zap.Error(err))
		return nil, err
	}
	return nil, nil
}
func (*CommonService) TagList(ctx context.Context, req *common.TagListReq) (*common.TagListResp, error) {

	var tags []*common.Tag
	user, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB

	if req.Name != "" {
		db = db.Where(`name LIKE ?` + "%" + req.Name + "%")
	}
	if req.GroupId != 0 {
		db = db.Where(`id IN ?`, db.Table(model.TableNameTagGroup).Select("tag_id").Where("group_id = ?", req.GroupId))
	}
	var count int64
	err = db.Table(`tag`).Where("user_id = ?", user.Id).Find(&tags).Count(&count).Error
	if err != nil {
		log.Errorw("db.Find", zap.Error(err))
		return nil, err
	}
	return &common.TagListResp{List: tags, Total: uint32(count)}, nil
}

func (*CommonService) SendMail(context.Context, *common.SendMailReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMail not implemented")
}
