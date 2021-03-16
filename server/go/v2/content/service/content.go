package service

import (
	"context"
	"net/http"

	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/protobuf/content"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContentService struct {
	content.UnimplementedContentServiceServer
}

func (m *ContentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "内容相关", "/api/content", nil
}


func (*ContentService) TagInfo(context.Context, *content.GetTagReq) (*content.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (*ContentService) AddTag(ctx context.Context, req *content.AddTagReq) (*empty.Empty, error) {
	ctxi, span := model.CtxFromContext(ctx).StartSpan("Edit")
	defer span.End()
	ctx = ctxi.Context
	user, err := ctxi.GetAuthInfo(Auth)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GORMDB
	req.UserId = user.Id
	err = db.Create(req).Error
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "db.Create", err.Error())
	}
	return nil, nil
}
func (*ContentService) EditTag(ctx context.Context, req *content.EditTagReq) (*empty.Empty, error) {
	ctxi, span := model.CtxFromContext(ctx).StartSpan("")
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
	}).Where(`id = ? AND user_id = ? AND status = 0`, req.Id, user.Id).Error
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "db.Updates", err.Error())
	}
	return nil, nil
}
func (*ContentService) TagList(ctx context.Context, req *content.TagListReq) (*content.TagListRep, error) {
	ctxi := model.CtxFromContext(ctx)
	var tags []*content.Tag

	user, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GORMDB

	if req.Name != "" {
		db = db.Where(`name LIKE ?` + "%" + req.Name + "%")
	}
	if req.Type != content.TagPlaceholder {
		db = db.Where(`type = ?`, req.Type)
	}
	var count int64
	err = db.Table(`tag`).Where("user_id = ?",user.Id).Find(&tags).Count(&count).Error
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "db.Find", err.Error())
	}
	return &content.TagListRep{List: tags, Total: uint32(count)}, nil
}
