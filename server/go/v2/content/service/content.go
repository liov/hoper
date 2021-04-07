package service

import (
	"context"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"net/http"

	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
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

func (*ContentService) TagInfo(context.Context, *request.Object) (*content.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (*ContentService) AddTag(ctx context.Context, req *content.AddTagReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("Edit")
	defer span.End()
	ctx = ctxi.Context
	user, err := auth(ctxi,false)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GORMDB
	req.UserId = user.Id
	err = db.Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "db.Create")
	}
	return nil, nil
}
func (*ContentService) EditTag(ctx context.Context, req *content.EditTagReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	ctx = ctxi.Context
	db := dao.Dao.GORMDB
	err = db.Updates(&content.Tag{
		Description:   req.Description,
		ExpressionURL: req.ExpressionURL,
	}).Where(`id = ? AND user_id = ? AND status = 0`, req.Id, auth.Id).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "db.Updates")
	}
	return nil, nil
}
func (*ContentService) TagList(ctx context.Context, req *content.TagListReq) (*content.TagListRep, error) {
	ctxi := contexti.CtxFromContext(ctx)
	var tags []*content.Tag

	user, err := auth(ctxi,true)
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
	err = db.Table(`tag`).Where("user_id = ?", user.Id).Find(&tags).Count(&count).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "db.Find")
	}
	return &content.TagListRep{List: tags, Total: uint32(count)}, nil
}

func (*ContentService) AddFav(ctx context.Context, req *content.AddFavReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	req.UserId = auth.Id
	err = db.Table(model.FavoritesTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
	}
	return nil, nil
}
func (*ContentService) EditFav(ctx context.Context, req *content.AddFavReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	err = db.Table(model.FavoritesTableName).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "UpdateColumn")
	}
	return nil, nil
}

// 创建合集
func (*ContentService) AddContainer(ctx context.Context,req *content.AddContainerReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	req.UserId = auth.Id
	err = db.Table(model.ContainerTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
	}
	return nil,nil
}

// 修改日记本
func (*ContentService) EditDiaryContainer(ctx context.Context,req *content.AddContainerReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	err = db.Table(model.FavoritesTableName).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
	}
	return nil,nil
}