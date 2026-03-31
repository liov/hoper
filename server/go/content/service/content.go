package service

import (
	"context"
	"net/http"

	"github.com/hopeio/gox/log"
	"github.com/hopeio/protobuf/request"
	"github.com/hopeio/scaffold/errcode"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

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

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	contentDBDao := data.GetDBDao(db)

	req.UserId = auth.Id
	id, err := contentDBDao.FavExists(ctx, req.Title, auth.Id)
	if err != nil {
		return nil, err
	}
	if id != 0 {
		return &request.Id{Id: id}, nil
	}

	err = contentDBDao.Transaction(func(tx *gorm.DB) error {
		contenttxDBDao := data.GetDBDao(tx)
		err = tx.Table(model.TableNameFavorite).Create(req).Error
		if err != nil {
			log.Errorw("Create", zap.Error(err))
			return errcode.DBError.Wrap(err)
		}
		err = contenttxDBDao.CreateContextExt(ctx, content.ContentFavorites, req.Id)
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
	//metadata := contextx.GetMetadata[*userpb.AuthInfo](ctx)

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(global.Dao.Redis.Client)
	err = contentRedisDao.Limit(ctx, &global.Conf.Moment.Limit, auth.Id)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	err = db.Table(model.TableNameFavorite).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		log.Errorw("Create", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}

// 收藏夹列表
func (*ContentService) FavList(ctx context.Context, req *content.FavListReq) (*content.FavListResp, error) {
	return nil, nil
}

// 收藏夹列表
func (*ContentService) TinyFavList(ctx context.Context, req *content.FavListReq) (*content.TinyFavListResp, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	var favs []*content.TinyFavorites
	if req.UserId == 0 {
		err = db.Table(model.TableNameFavorite).Select("id,title").Where(`user_id = ?`, auth.Id).Find(&favs).Error
	}
	if err != nil {
		log.Errorw("Find", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	return &content.TinyFavListResp{List: favs}, nil
}

// 创建合集
func (*ContentService) AddSet(ctx context.Context, req *content.AddSetReq) (*emptypb.Empty, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(global.Dao.Redis.Client)
	err = contentRedisDao.Limit(ctx, &global.Conf.Moment.Limit, auth.Id)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	req.UserId = auth.Id
	err = db.Table(model.TableNameContainer).Create(req).Error
	if err != nil {
		log.Errorw("Create", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}

// 修改合集
func (*ContentService) EditSet(ctx context.Context, req *content.AddSetReq) (*emptypb.Empty, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(global.Dao.Redis.Client)
	err = contentRedisDao.Limit(ctx, &global.Conf.Moment.Limit, auth.Id)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	err = db.Table(model.TableNameContainer).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		log.Errorw("Create", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}
