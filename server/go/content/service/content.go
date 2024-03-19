package service

import (
	"context"
	"github.com/hopeio/tiga/context/http_context"
	"github.com/hopeio/tiga/protobuf/request"
	gormi "github.com/hopeio/tiga/utils/dao/db/gorm"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"net/http"

	"github.com/liov/hoper/server/go/content/confdao"
	"github.com/liov/hoper/server/go/content/data"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"

	"github.com/hopeio/tiga/protobuf/errorcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContentService struct {
	content.UnimplementedContentServiceServer
}

func (m *ContentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "内容相关", "/api/content", nil
}

func (*ContentService) TagInfo(context.Context, *request.Id) (*content.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (*ContentService) AddTag(ctx context.Context, req *content.AddTagReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("Edit")
	defer span.End()
	ctx = ctxi.Context
	user, err := auth(ctxi, false)
	if err != nil {
		return nil, err
	}
	db := confdao.Dao.GORMDB.DB
	req.UserId = user.Id
	err = db.Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "db.Create")
	}
	return nil, nil
}
func (*ContentService) EditTag(ctx context.Context, req *content.EditTagReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	ctx = ctxi.Context
	db := confdao.Dao.GORMDB.DB
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
	ctxi := http_context.ContextFromContext(ctx)
	var tags []*content.Tag

	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := confdao.Dao.GORMDB.DB

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

func (*ContentService) AddFav(ctx context.Context, req *content.AddFavReq) (*request.Id, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	contentDBDao := data.GetDBDao(ctxi, db)

	req.UserId = auth.Id
	id, err := contentDBDao.FavExists(req.Title)
	if err != nil {
		return nil, err
	}
	if id != 0 {
		return &request.Id{Id: id}, nil
	}

	err = contentDBDao.Transaction(func(tx *gorm.DB) error {
		contenttxDBDao := data.GetDBDao(ctxi, tx)
		err = tx.Table(model.FavoritesTableName).Create(req).Error
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
		}
		err = contenttxDBDao.CreateContextExt(content.ContentFavorites, req.Id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &request.Id{Id: req.Id}, nil
}
func (*ContentService) EditFav(ctx context.Context, req *content.AddFavReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, confdao.Dao.Redis)
	err = contentRedisDao.Limit(&confdao.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	err = db.Table(model.FavoritesTableName).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "UpdateColumn")
	}
	return nil, nil
}

// 收藏夹列表
func (*ContentService) FavList(ctx context.Context, req *content.FavListReq) (*content.FavListRep, error) {
	return nil, nil
}

// 收藏夹列表
func (*ContentService) TinyFavList(ctx context.Context, req *content.FavListReq) (*content.TinyFavListRep, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	var favs []*content.TinyFavorites
	if req.UserId == 0 {
		err = db.Table(model.FavoritesTableName).Select("id,title").Where(`user_id = ?`, auth.Id).Find(&favs).Error
	}
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
	}
	return &content.TinyFavListRep{List: favs}, nil
}

// 创建合集
func (*ContentService) AddContainer(ctx context.Context, req *content.AddContainerReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, confdao.Dao.Redis)
	err = contentRedisDao.Limit(&confdao.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	req.UserId = auth.Id
	err = db.Table(model.ContainerTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
	}
	return nil, nil
}

// 修改合集
func (*ContentService) EditDiaryContainer(ctx context.Context, req *content.AddContainerReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, confdao.Dao.Redis)
	err = contentRedisDao.Limit(&confdao.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	err = db.Table(model.ContainerTableName).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "CreateFav")
	}
	return nil, nil
}
