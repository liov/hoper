package service

import (
	"context"

	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/protobuf/content"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContentService struct {
	content.UnimplementedContentServiceServer
}

func (*ContentService) Info(context.Context, *content.GetTagReq) (*content.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (*ContentService) Add(ctx context.Context, req *content.AddTagReq) (*request.Empty, error) {
	ctxi, span := model.CtxFromContext(ctx).StartSpan("Edit")
	defer span.End()
	ctx = ctxi.Context
	user, err := ctxi.GetAuthInfo(Auth)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GORMDB
	db.Create(&content.Tag{
		Name:          req.Name,
		Description:   req.Description,
		ExpressionURL: req.ExpressionURL,
		UserId:        user.Id,
	})
	return nil, nil
}
func (*ContentService) Edit(ctx context.Context, req *content.AddTagReq) (*request.Empty, error) {
	ctxi, span := model.CtxFromContext(ctx).StartSpan("Edit")
	defer span.End()
	user, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	ctx = ctxi.Context
	db := dao.Dao.GORMDB
	err = db.Updates(&content.Tag{
		Description:   req.Description,
		ExpressionURL: req.ExpressionURL,
	}).Where(`name = ? AND user_id = ? AND status = 0`, req.Name, user.Id).Error
	if err != nil {
		return nil, errorcode.DBError.Warp(err)
	}
	return nil, nil
}
func (*ContentService) List(ctx context.Context, req *content.TagListReq) (*content.TagListRep, error) {
	var tags []*content.Tag
	db := dao.Dao.GORMDB

	if req.Name != "" {
		db = db.Where(`name LIKE ?` + "%" + req.Name + "%")
	}
	var count int64
	db.Table(`tag`).Find(&tags).Count(&count)
	return &content.TagListRep{List: tags, Count: uint32(count)}, nil
}
