package service

import (
	"context"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/protobuf/request"
	"github.com/hopeio/scaffold/errcode"
	gormi "github.com/hopeio/utils/datax/database/gorm"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"net/http"

	"github.com/liov/hoper/server/go/content/data"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/content"
)

type ContentService struct {
	content.UnimplementedContentServiceServer
}

func (m *ContentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "内容相关", "/api/content", nil
}

func (*ContentService) AddFav(ctx context.Context, req *content.AddFavReq) (*request.Id, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
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
		err = tx.Table(model.TableNameFavorite).Create(req).Error
		if err != nil {
			return ctxi.RespErrorLog(errcode.DBError, err, "CreateFav")
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
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.Limit(&global.Conf.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	err = db.Table(model.TableNameFavorite).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "UpdateColumn")
	}
	return nil, nil
}

// 收藏夹列表
func (*ContentService) FavList(ctx context.Context, req *content.FavListReq) (*content.FavListRep, error) {
	return nil, nil
}

// 收藏夹列表
func (*ContentService) TinyFavList(ctx context.Context, req *content.FavListReq) (*content.TinyFavListRep, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	var favs []*content.TinyFavorites
	if req.UserId == 0 {
		err = db.Table(model.TableNameFavorite).Select("id,title").Where(`user_id = ?`, auth.Id).Find(&favs).Error
	}
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "CreateFav")
	}
	return &content.TinyFavListRep{List: favs}, nil
}

// 创建合集
func (*ContentService) AddSet(ctx context.Context, req *content.AddSetReq) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.Limit(&global.Conf.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	req.UserId = auth.Id
	err = db.Table(model.TableNameContainer).Create(req).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "CreateFav")
	}
	return nil, nil
}

// 修改合集
func (*ContentService) EditSet(ctx context.Context, req *content.AddSetReq) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.Limit(&global.Conf.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	err = db.Table(model.TableNameContainer).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "CreateFav")
	}
	return nil, nil
}
